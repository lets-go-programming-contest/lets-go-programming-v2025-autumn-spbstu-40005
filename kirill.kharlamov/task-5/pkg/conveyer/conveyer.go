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
	ErrNoDataAvailable   = errors.New("no data available")
)

type Conveyer struct {
	bufferSize int
	channels   map[string]chan string
	tasks      []func(context.Context) error
	mutex      sync.RWMutex
}

func New(size int) *Conveyer {
	return &Conveyer{
		bufferSize: size,
		channels:   make(map[string]chan string),
		tasks:      make([]func(context.Context) error, 0),
		mutex:      sync.RWMutex{},
	}
}

func (c *Conveyer) RegisterDecorator(
	decoratorFunc func(context.Context, chan string, chan string) error,
	inputName, outputName string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inputChannel := c.ensureChannel(inputName)
	outputChannel := c.ensureChannel(outputName)

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return decoratorFunc(ctx, inputChannel, outputChannel)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	multiplexerFunc func(context.Context, []chan string, chan string) error,
	inputNames []string,
	outputName string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	outputChannel := c.ensureChannel(outputName)
	inputChannels := make([]chan string, len(inputNames))

	for i, name := range inputNames {
		inputChannels[i] = c.ensureChannel(name)
	}

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return multiplexerFunc(ctx, inputChannels, outputChannel)
	})
}

func (c *Conveyer) RegisterSeparator(
	separatorFunc func(context.Context, chan string, []chan string) error,
	inputName string,
	outputNames []string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inputChannel := c.ensureChannel(inputName)
	outputChannels := make([]chan string, len(outputNames))

	for i, name := range outputNames {
		outputChannels[i] = c.ensureChannel(name)
	}

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return separatorFunc(ctx, inputChannel, outputChannels)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	errGroup, ctx := errgroup.WithContext(ctx)

	for _, task := range c.tasks {
		taskFunc := task

		errGroup.Go(func() error {
			return taskFunc(ctx)
		})
	}

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

	for channelName, channel := range c.channels {
		close(channel)
		delete(c.channels, channelName)
	}
}

func (c *Conveyer) Send(channelName string, data string) error {
	c.mutex.RLock()
	channel, exists := c.channels[channelName]
	c.mutex.RUnlock()

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
	c.mutex.RLock()
	channel, exists := c.channels[channelName]
	c.mutex.RUnlock()

	if !exists {
		return "", ErrChannelNotFound
	}

	select {
	case value, ok := <-channel:
		if !ok {
			return emptyValue, nil
		}

		return value, nil
	default:
		return "", ErrNoDataAvailable
	}
}

func (c *Conveyer) ensureChannel(name string) chan string {
	if channel, exists := c.channels[name]; exists {
		return channel
	}

	newChannel := make(chan string, c.bufferSize)
	c.channels[name] = newChannel

	return newChannel
}
