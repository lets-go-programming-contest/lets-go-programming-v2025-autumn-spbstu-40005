package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const undefined = "undefined"

var ErrFoundOfChannel = errors.New("chan not found")

type conveyerImpl struct {
	size     int
	channels map[string]chan string
	handlers []func(ctx context.Context) error
	mu       sync.RWMutex
}

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		size:     size,
		channels: make(map[string]chan string),
		handlers: []func(ctx context.Context) error{},
		mu:       sync.RWMutex{},
	}
}

func (c *conveyerImpl) RegisterDecorator(
	functional func(ctx context.Context, input chan string, output chan string) error,
	input,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputCh := c.getOrCreateChannel(input)
	outputCh := c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return functional(ctx, inputCh, outputCh)
	})
}

func (c *conveyerImpl) RegisterMultiplexer(
	functional func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	outputCh := c.getOrCreateChannel(output)
	inputChs := make([]chan string, len(inputs))

	for i, name := range inputs {
		inputChs[i] = c.getOrCreateChannel(name)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return functional(ctx, inputChs, outputCh)
	})
}

func (c *conveyerImpl) RegisterSeparator(
	functional func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputCh := c.getOrCreateChannel(input)
	outputChs := make([]chan string, len(outputs))

	for i, name := range outputs {
		outputChs[i] = c.getOrCreateChannel(name)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return functional(ctx, inputCh, outputChs)
	})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	defer c.closeAllChannels()

	errgr, ctx := errgroup.WithContext(ctx)

	for _, h := range c.handlers {
		errgr.Go(func() error {
			return h(ctx)
		})
	}

	err := errgr.Wait()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (c *conveyerImpl) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, channel := range c.channels {
		close(channel)
	}
}

func (c *conveyerImpl) Send(input string, data string) error {
	channel, exists := c.getChannelForRead(input)

	if !exists {
		return ErrFoundOfChannel
	}

	channel <- data

	return nil
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	channel, exists := c.getChannelForRead(output)

	if !exists {
		return "", ErrFoundOfChannel
	}

	val, ok := <-channel
	if !ok {
		return undefined, nil
	}

	return val, nil
}

func (c *conveyerImpl) getOrCreateChannel(name string) chan string {
	if channel, ok := c.channels[name]; ok {
		return channel
	}

	channel := make(chan string, c.size)
	c.channels[name] = channel

	return channel
}

func (c *conveyerImpl) getChannelForRead(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ch, ok := c.channels[name]
	return ch, ok
}
