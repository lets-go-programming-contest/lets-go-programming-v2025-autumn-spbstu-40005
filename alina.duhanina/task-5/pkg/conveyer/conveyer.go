package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound   = errors.New("chan not found")
	ErrChannelFull    = errors.New("channel is full")
	ErrNoDataAvailable = errors.New("no data available")
	ErrUndefined      = "undefined"
)

type DecoratorFunc func(ctx context.Context, input chan string, output chan string) error
type MultiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error
type SeparatorFunc func(ctx context.Context, input chan string, outputs []chan string) error

type Conveyer struct {
	size     int
	channels map[string]chan string
	mu       sync.RWMutex
	wg       sync.WaitGroup
	cancel   context.CancelFunc
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
	}
}

func (c *Conveyer) getOrCreateChannel(id string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[id]; exists {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[id] = ch
	return ch
}

func (c *Conveyer) getChannel(id string) (chan string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if ch, exists := c.channels[id]; exists {
		return ch, nil
	}

	return nil, ErrChanNotFound
}

func (c *Conveyer) RegisterDecorator(
	fn DecoratorFunc,
	input string,
	output string,
) {
	c.getOrCreateChannel(input)
	c.getOrCreateChannel(output)
}

func (c *Conveyer) RegisterMultiplexer(
	fn MultiplexerFunc,
	inputs []string,
	output string,
) {
	for _, input := range inputs {
		c.getOrCreateChannel(input)
	}
	c.getOrCreateChannel(output)
}

func (c *Conveyer) RegisterSeparator(
	fn SeparatorFunc,
	input string,
	outputs []string,
) {
	c.getOrCreateChannel(input)
	for _, output := range outputs {
		c.getOrCreateChannel(output)
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, c.cancel = context.WithCancel(ctx)
	defer c.cancel()

	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Conveyer) Send(input string, data string) error {
	ch, err := c.getChannel(input)
	if err != nil {
		return ErrChanNotFound
	}

	select {
	case ch <- data:
		return nil
	default:
		return ErrChannelFull
	}
}

func (c *Conveyer) Recv(output string) (string, error) {
	ch, err := c.getChannel(output)
	if err != nil {
		return "", ErrChanNotFound
	}

	select {
	case data, ok := <-ch:
		if !ok {
			return ErrUndefined, nil
		}
		return data, nil
	default:
		return "", ErrNoDataAvailable
	}
}

func (c *Conveyer) Stop() {
	if c.cancel != nil {
		c.cancel()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for id, ch := range c.channels {
		close(ch)
		delete(c.channels, id)
	}
}
