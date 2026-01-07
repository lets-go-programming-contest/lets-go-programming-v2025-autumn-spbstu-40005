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

	mu       sync.RWMutex
	channels map[string]chan string
	handlers []func(context.Context) error
}

func New(size int) Conveyer {
	if size < 0 {
		size = 0
	}

	return Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]func(context.Context) error, 0),
	}
}

func (c *Conveyer) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.makeChannels(input, output)

	inCh := c.channels[input]
	outCh := c.channels[output]

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.makeChannels(output)
	c.makeChannels(inputs...)

	inChans := make([]chan string, 0, len(inputs))
	for _, name := range inputs {
		inChans = append(inChans, c.channels[name])
	}

	outCh := c.channels[output]

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inChans, outCh)
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.makeChannels(input)
	c.makeChannels(outputs...)

	inCh := c.channels[input]

	outChans := make([]chan string, 0, len(outputs))
	for _, name := range outputs {
		outChans = append(outChans, c.channels[name])
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inCh, outChans)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	defer c.closeAllChannels()

	group, runCtx := errgroup.WithContext(ctx)

	c.mu.RLock()
	handlers := make([]func(context.Context) error, len(c.handlers))
	copy(handlers, c.handlers)
	c.mu.RUnlock()

	for _, h := range handlers {
		handler := h
		group.Go(func() error {
			return handler(runCtx)
		})
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("conveyer run: %w", err)
	}

	return nil
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.RLock()
	ch, ok := c.channels[input]
	c.mu.RUnlock()

	if !ok {
		return ErrNoChannel
	}

	ch <- data
	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, ok := c.channels[output]
	c.mu.RUnlock()

	if !ok {
		return "", ErrNoChannel
	}

	data, opened := <-ch
	if !opened {
		return undefined, nil
	}

	return data, nil
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

	for _, ch := range c.channels {
		close(ch)
	}
}
