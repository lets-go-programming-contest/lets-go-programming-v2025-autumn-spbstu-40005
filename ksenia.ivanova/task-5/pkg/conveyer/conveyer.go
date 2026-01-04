package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrChanNotFound   = errors.New("chan not found")
	ErrUndefined      = errors.New("undefined")
	ErrAlreadyStarted = errors.New("conveyer already started")
)

type Conveyer struct {
	mu         sync.RWMutex
	size       int
	channels   map[string]chan string
	processors []func(context.Context) error
	started    bool
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
	}
}

func (c *Conveyer) getOrCreateChan(id string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[id]; exists {
		return ch
	}
	ch := make(chan string, c.size)
	c.channels[id] = ch
	return ch
}

func (c *Conveyer) getChan(id string) (chan string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, exists := c.channels[id]
	if !exists {
		return nil, ErrChanNotFound
	}
	return ch, nil
}

func (c *Conveyer) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	inputID, outputID string,
) {
	c.registerProcessor(func(ctx context.Context) error {
		input := c.getOrCreateChan(inputID)
		output := c.getOrCreateChan(outputID)
		return fn(ctx, input, output)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputIDs []string, outputID string,
) {
	c.registerProcessor(func(ctx context.Context) error {
		inputs := make([]chan string, len(inputIDs))
		for i, id := range inputIDs {
			inputs[i] = c.getOrCreateChan(id)
		}
		output := c.getOrCreateChan(outputID)
		return fn(ctx, inputs, output)
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	inputID string, outputIDs []string,
) {
	c.registerProcessor(func(ctx context.Context) error {
		input := c.getOrCreateChan(inputID)
		outputs := make([]chan string, len(outputIDs))
		for i, id := range outputIDs {
			outputs[i] = c.getOrCreateChan(id)
		}
		return fn(ctx, input, outputs)
	})
}

func (c *Conveyer) registerProcessor(fn func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.started {
		return
	}
	c.processors = append(c.processors, fn)
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.started {
		c.mu.Unlock()
		return ErrAlreadyStarted
	}
	c.started = true
	c.ctx, c.cancel = context.WithCancel(ctx)
	c.mu.Unlock()

	defer c.cleanup()

	errCh := make(chan error, len(c.processors))

	for _, proc := range c.processors {
		c.wg.Add(1)
		go func(p func(context.Context) error) {
			defer c.wg.Done()
			if err := p(c.ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
				c.cancel()
			}
		}(proc)
	}

	select {
	case err := <-errCh:
		return err
	case <-c.ctx.Done():
		return c.ctx.Err()
	}
}

func (c *Conveyer) cleanup() {
	c.wg.Wait()
	c.mu.Lock()
	defer c.mu.Unlock()
	for id, ch := range c.channels {
		close(ch)
		delete(c.channels, id)
	}
}

func (c *Conveyer) Send(inputID, data string) error {
	ch, err := c.getChan(inputID)
	if err != nil {
		return err
	}

	select {
	case ch <- data:
		return nil
	case <-c.ctx.Done():
		return c.ctx.Err()
	}
}

func (c *Conveyer) Recv(outputID string) (string, error) {
	ch, err := c.getChan(outputID)
	if err != nil {
		return "", err
	}

	select {
	case data, ok := <-ch:
		if !ok {
			return "undefined", nil
		}
		return data, nil
	case <-c.ctx.Done():
		return "", c.ctx.Err()
	}
}
