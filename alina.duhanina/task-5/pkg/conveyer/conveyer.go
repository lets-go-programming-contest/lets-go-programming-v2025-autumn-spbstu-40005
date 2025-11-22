package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound    = errors.New("chan not found")
	ErrChannelFull     = errors.New("channel is full")
	ErrNoDataAvailable = errors.New("no data available")
	ErrUndefined       = "undefined"
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
	inputChan := c.getOrCreateChannel(input)
	outputChan := c.getOrCreateChannel(output)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		_ = fn(context.Background(), inputChan, outputChan)
	}()
}

func (c *Conveyer) RegisterMultiplexer(
	fn MultiplexerFunc,
	inputs []string,
	output string,
) {
	inputChans := make([]chan string, len(inputs))
	for i, input := range inputs {
		inputChans[i] = c.getOrCreateChannel(input)
	}
	outputChan := c.getOrCreateChannel(output)

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		_ = fn(context.Background(), inputChans, outputChan)
	}()
}

func (c *Conveyer) RegisterSeparator(
	fn SeparatorFunc,
	input string,
	outputs []string,
) {
	inputChan := c.getOrCreateChannel(input)
	outputChans := make([]chan string, len(outputs))
	for i, output := range outputs {
		outputChans[i] = c.getOrCreateChannel(output)
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		_ = fn(context.Background(), inputChan, outputChans)
	}()
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, c.cancel = context.WithCancel(ctx)
	defer c.cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		c.wg.Wait()
		return nil
	})

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-func() chan struct{} {
		done := make(chan struct{})
		go func() {
			g.Wait()
			close(done)
		}()
		return done
	}():
		return nil
	}
}

func (c *Conveyer) RunParallel(ctx context.Context, handlers []func(context.Context) error) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, handler := range handlers {
		handler := handler
		g.Go(func() error {
			return handler(ctx)
		})
	}

	return g.Wait()
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
