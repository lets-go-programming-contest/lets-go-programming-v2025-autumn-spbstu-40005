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

var ErrChannelNotFound = errors.New("channel not found")

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
		errGroup.Go(func() error {
			return task(ctx)
		})
	}

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("conveyer execution failed: %w", err)
	}

	c.closeAllChannels()

	return nil
}

func (c *Conveyer) closeAllChannels() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, channel := range c.channels {
		close(channel)
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
	case <-context.TODO().Done():
		return context.Canceled
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
	case val, ok := <-channel:
		if !ok {
			return emptyValue, nil
		}
		return val, nil
	case <-context.TODO().Done():
		return "", context.Canceled
	}
}

func (c *Conveyer) ensureChannel(name string) chan string {
	c.mutex.RLock()
	if ch, exists := c.channels[name]; exists {
		c.mutex.RUnlock()
		return ch
	}
	c.mutex.RUnlock()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.bufferSize)
	c.channels[name] = ch
	return ch
}
