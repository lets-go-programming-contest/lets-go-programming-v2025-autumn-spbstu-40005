package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound = errors.New("chan not found")
	ErrUndefined    = "undefined"
)

type Conveyer struct {
	mu       sync.RWMutex
	channels map[string]chan string
	size     int
	handlers []func(ctx context.Context) error
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		size:     size,
	}
}

func (c *Conveyer) getOrCreateChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}
	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) getChan(name string) (chan string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, exists := c.channels[name]
	if !exists {
		return nil, ErrChanNotFound
	}
	return ch, nil
}

func (c *Conveyer) RegisterDecorator(fn func(ctx context.Context, input, output chan string) error, input, output string) {
	in := c.getOrCreateChan(input)
	out := c.getOrCreateChan(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *Conveyer) RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string) {
	inChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChans[i] = c.getOrCreateChan(name)
	}
	out := c.getOrCreateChan(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inChans, out)
	})
}

func (c *Conveyer) RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string) {
	in := c.getOrCreateChan(input)
	outChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChans[i] = c.getOrCreateChan(name)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, in, outChans)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, handler := range c.handlers {
		h := handler
		g.Go(func() error {
			return h(ctx)
		})
	}

	return g.Wait()
}

func (c *Conveyer) Send(input string, data string) error {
	ch, err := c.getChan(input)
	if err != nil {
		return err
	}

	select {
	case ch <- data:
		return nil
	default:
		return errors.New("channel is full")
	}
}

func (c *Conveyer) Recv(output string) (string, error) {
	ch, err := c.getChan(output)
	if err != nil {
		return "", err
	}

	select {
	case data, ok := <-ch:
		if !ok {
			return ErrUndefined, nil
		}
		return data, nil
	default:
		return "", errors.New("no data available")
	}
}
