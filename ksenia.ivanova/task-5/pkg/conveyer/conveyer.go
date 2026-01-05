package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrChanNotFound = errors.New("chan not found")
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
	errChan    chan error
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		errChan:  make(chan error, 1),
	}
}

func (c *Conveyer) getOrCreateChan(id string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.channels[id]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[id] = ch
	return ch
}

func (c *Conveyer) getChan(id string) (chan string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, ok := c.channels[id]
	if !ok {
		return nil, ErrChanNotFound
	}
	return ch, nil
}

func (c *Conveyer) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	inputID, outputID string,
) {
	c.getOrCreateChan(inputID)
	c.getOrCreateChan(outputID)

	c.processors = append(c.processors, func(ctx context.Context) error {
		input := c.getOrCreateChan(inputID)
		output := c.getOrCreateChan(outputID)
		return fn(ctx, input, output)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputIDs []string, outputID string,
) {
	for _, id := range inputIDs {
		c.getOrCreateChan(id)
	}
	c.getOrCreateChan(outputID)

	c.processors = append(c.processors, func(ctx context.Context) error {
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
	c.getOrCreateChan(inputID)
	for _, id := range outputIDs {
		c.getOrCreateChan(id)
	}

	c.processors = append(c.processors, func(ctx context.Context) error {
		input := c.getOrCreateChan(inputID)
		outputs := make([]chan string, len(outputIDs))
		for i, id := range outputIDs {
			outputs[i] = c.getOrCreateChan(id)
		}
		return fn(ctx, input, outputs)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.started {
		c.mu.Unlock()
		return errors.New("conveyer already started")
	}
	c.started = true
	c.ctx, c.cancel = context.WithCancel(ctx)
	c.mu.Unlock()

	done := make(chan struct{})

	for _, proc := range c.processors {
		c.wg.Add(1)
		go func(p func(context.Context) error) {
			defer c.wg.Done()
			if err := p(c.ctx); err != nil {
				select {
				case c.errChan <- err:
				default:
				}
				c.cancel()
			}
		}(proc)
	}

	go func() {
		c.wg.Wait()
		close(done)
	}()

	var err error
	select {
	case err = <-c.errChan:
		c.cancel()
	case <-done:
		err = nil
	}

	c.closeChannels()
	return err
}

func (c *Conveyer) closeChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, ch := range c.channels {
		close(ch)
	}
}

func (c *Conveyer) Send(inputID, data string) error {
	ch, err := c.getChan(inputID)
	if err != nil {
		return err
	}

	c.mu.RLock()
	started := c.started
	ctx := c.ctx
	c.mu.RUnlock()

	if !started {
		select {
		case ch <- data:
			return nil
		default:
			return errors.New("send failed")
		}
	}

	select {
	case ch <- data:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Conveyer) Recv(outputID string) (string, error) {
	ch, err := c.getChan(outputID)
	if err != nil {
		return "", err
	}

	c.mu.RLock()
	started := c.started
	ctx := c.ctx
	c.mu.RUnlock()
	if !started {
		select {
		case data, ok := <-ch:
			if !ok {
				return "", errors.New("channel closed")
			}
			return data, nil
		default:
			return "", errors.New("no data")
		}
	}

	select {
	case data, ok := <-ch:
		if !ok {
			return "", errors.New("channel closed")
		}
		return data, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
