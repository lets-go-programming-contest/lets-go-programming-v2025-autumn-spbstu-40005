package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrNoChannel = errors.New("chan not found")

const undefined = "undefined"

type Conveyer struct {
	size int

	channels map[string]chan string
	handlers []func(context.Context) error
	mu       sync.RWMutex
}

func New(size int) Conveyer {
	if size < 0 {
		size = 0
	}

	return Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]func(context.Context) error, 0),
		mu:       sync.RWMutex{},
	}
}

func (c *Conveyer) RegisterDecorator(
	handlerFunc func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.makeChannels(input, output)

	inputChannel := c.channels[input]
	outputChannel := c.channels[output]

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannel, outputChannel)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	handlerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.makeChannels(output)
	c.makeChannels(inputs...)

	inputChannels := make([]chan string, 0, len(inputs))
	for _, name := range inputs {
		inputChannels = append(inputChannels, c.channels[name])
	}

	outputChannel := c.channels[output]

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannels, outputChannel)
	})
}

func (c *Conveyer) RegisterSeparator(
	handlerFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.makeChannels(input)
	c.makeChannels(outputs...)

	inputChannel := c.channels[input]

	outputChannels := make([]chan string, 0, len(outputs))
	for _, name := range outputs {
		outputChannels = append(outputChannels, c.channels[name])
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannel, outputChannels)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	defer c.closeAllChannels()

	group, groupCtx := errgroup.WithContext(ctx)

	handlers := c.snapshotHandlers()
	for _, handlerFunc := range handlers {
		currentHandler := handlerFunc

		task := func() error {
			return currentHandler(groupCtx)
		}

		group.Go(task)
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("conveyer run: %w", err)
	}

	return nil
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.RLock()
	channel, ok := c.channels[input]
	c.mu.RUnlock()

	if !ok {
		return ErrNoChannel
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	channel, ok := c.channels[output]
	c.mu.RUnlock()

	if !ok {
		return "", ErrNoChannel
	}

	data, opened := <-channel
	if !opened {
		return undefined, nil
	}

	return data, nil
}

func (c *Conveyer) snapshotHandlers() []func(context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	handlers := make([]func(context.Context) error, len(c.handlers))
	copy(handlers, c.handlers)

	return handlers
}

func (c *Conveyer) makeChannel(name string) {
	if _, ok := c.channels[name]; ok {
		return
	}

	c.channels[name] = make(chan string, c.size)
}

func (c *Conveyer) makeChannels(names ...string) {
	for _, name := range names {
		c.makeChannel(name)
	}
}

func (c *Conveyer) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, channel := range c.channels {
		close(channel)
	}
}
