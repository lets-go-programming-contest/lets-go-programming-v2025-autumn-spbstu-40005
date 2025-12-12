package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const (
	emptyValue = "undefined"
)

var (
	ErrChannelNotFound   = errors.New("chan not found")
	ErrChannelBufferFull = errors.New("channel buffer is full")
)

type taskFunc func(context.Context) error

type Conveyer struct {
	bufferSize int
	channels   map[string]chan string
	tasks      []taskFunc
	mutex      sync.RWMutex
}

func New(size int) *Conveyer {
	return &Conveyer{
		bufferSize: size,
		channels:   make(map[string]chan string),
		tasks:      make([]taskFunc, 0),
	}
}

func (c *Conveyer) RegisterDecorator(
	decoratorFunc func(context.Context, chan string, chan string) error,
	inputName, outputName string,
) {
	inputChannel := c.ensureChannel(inputName)
	outputChannel := c.ensureChannel(outputName)

	c.mutex.Lock()
	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return decoratorFunc(ctx, inputChannel, outputChannel)
	})
	c.mutex.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(
	multiplexerFunc func(context.Context, []chan string, chan string) error,
	inputNames []string,
	outputName string,
) {
	outputChannel := c.ensureChannel(outputName)
	inputChannels := make([]chan string, len(inputNames))

	for i, name := range inputNames {
		inputChannels[i] = c.ensureChannel(name)
	}

	c.mutex.Lock()
	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return multiplexerFunc(ctx, inputChannels, outputChannel)
	})
	c.mutex.Unlock()
}

func (c *Conveyer) RegisterSeparator(
	separatorFunc func(context.Context, chan string, []chan string) error,
	inputName string,
	outputNames []string,
) {
	inputChannel := c.ensureChannel(inputName)
	outputChannels := make([]chan string, len(outputNames))

	for i, name := range outputNames {
		outputChannels[i] = c.ensureChannel(name)
	}

	c.mutex.Lock()
	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return separatorFunc(ctx, inputChannel, outputChannels)
	})
	c.mutex.Unlock()
}

func (c *Conveyer) Run(ctx context.Context) error {
	errGroup, ctx := errgroup.WithContext(ctx)

	c.mutex.RLock()
	for _, task := range c.tasks {
		taskFunc := task
		errGroup.Go(func() error {
			return taskFunc(ctx)
		})
	}
	c.mutex.RUnlock()

	err := errGroup.Wait()
	c.closeAllChannels()

	if err != nil {
		return fmt.Errorf("conveyer execution failed: %w", err)
	}

	return nil
}

func (c *Conveyer) closeAllChannels() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, channel := range c.channels {
		close(channel)
	}
}

func (c *Conveyer) getChannel(name string) (chan string, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	channel, exists := c.channels[name]
	return channel, exists
}

func (c *Conveyer) Send(channelName string, data string) error {
	channel, exists := c.getChannel(channelName)

	if !exists {
		return ErrChannelNotFound
	}

	select {
	case channel <- data:
		return nil
	default:
		return ErrChannelBufferFull
	}
}

func (c *Conveyer) Recv(channelName string) (string, error) {
	channel, exists := c.getChannel(channelName)

	if !exists {
		return "", ErrChannelNotFound
	}

	value, ok := <-channel
	if !ok {
		return emptyValue, nil
	}

	return value, nil
}

func (c *Conveyer) ensureChannel(name string) chan string {
	c.mutex.RLock()
	channel, exists := c.channels[name]
	c.mutex.RUnlock()

	if exists {
		return channel
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if channel, exists := c.channels[name]; exists {
		return channel
	}

	newChannel := make(chan string, c.bufferSize)
	c.channels[name] = newChannel

	return newChannel
}
