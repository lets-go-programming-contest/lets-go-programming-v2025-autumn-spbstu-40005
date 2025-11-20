package conveyer

import (
	"context"
	"errors"
	"sync"
)

type conveyer interface {
	RegisterDecorator(
		fn func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		fn func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		fn func(ctx context.Context, input chan string, outputs []chan string) error,
		input string,
		outputs []string,
	)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type conveyerImpl struct {
	size     int
	channels map[string]chan string
	handlers []func(ctx context.Context) error
	mu       sync.RWMutex
}

func New(size int) conveyer {
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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, 1)
	var wg sync.WaitGroup

	for _, h := range c.handlers {
		wg.Add(1)
		go func(fn func(context.Context) error) {
			defer wg.Done()
			if err := fn(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
				cancel()
			}
		}(h)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			c.closeAllChannels()
			return err
		}

		<-ctx.Done()
	case <-ctx.Done():
		wg.Wait()
	}

	c.closeAllChannels()
	return ctx.Err()
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
