package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

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
	}
}

func (c *conveyerImpl) RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error, input, output string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputCh := c.getOrCreateChannel(input)
	outputCh := c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputCh, outputCh)
	})
}

func (c *conveyerImpl) RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	outputCh := c.getOrCreateChannel(output)
	inputChs := make([]chan string, len(inputs))
	for i, name := range inputs {
		inputChs[i] = c.getOrCreateChannel(name)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputChs, outputCh)
	})
}

func (c *conveyerImpl) RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputCh := c.getOrCreateChannel(input)
	outputChs := make([]chan string, len(outputs))
	for i, name := range outputs {
		outputChs[i] = c.getOrCreateChannel(name)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputCh, outputChs)
	})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	defer c.closeAllChannels()
	eg, ctx := errgroup.WithContext(ctx)
	for _, h := range c.handlers {
		eg.Go(func() error {
			return h(ctx)
		})
	}
	return eg.Wait()
}

func (c *conveyerImpl) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, ch := range c.channels {
		close(ch)
	}
}

func (c *conveyerImpl) Send(input string, data string) error {
	c.mu.RLock()
	ch, exists := c.channels[input]
	c.mu.RUnlock()

	if !exists {
		return errors.New("chan not found")
	}

	ch <- data
	return nil
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, exists := c.channels[output]
	c.mu.RUnlock()

	if !exists {
		return "", errors.New("chan not found")
	}

	val, ok := <-ch
	if !ok {
		return "undefined", nil
	}
	return val, nil
}

func (c *conveyerImpl) getOrCreateChannel(name string) chan string {
	if ch, ok := c.channels[name]; ok {
		return ch
	}
	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}
