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

	channelIds := make([]string, 0, len(c.channels))
	for id := range c.channels {
		channelIds = append(channelIds, id)
	}

	c.ctx, c.cancel = context.WithCancel(ctx)
	c.mu.Unlock()

	for _, id := range channelIds {
		c.getOrCreateChan(id)
	}

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

	var err error
	select {
	case err = <-c.errChan:
		c.cancel()
	case <-c.ctx.Done():
		err = c.ctx.Err()
	}

	c.wg.Wait()
	c.closeChannels()

	if errors.Is(err, context.Canceled) {
		return nil
	}
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
	if c.ctx == nil {
		ch := c.getOrCreateChan(inputID)
		select {
		case ch <- data:
			return nil
		default:
			return errors.New("send failed")
		}
	}

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
	if c.ctx == nil {
		ch := c.getOrCreateChan(outputID)
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

	ch, err := c.getChan(outputID)
	if err != nil {
		return "", err
	}

	select {
	case data, ok := <-ch:
		if !ok {
			return "", nil
		}
		return data, nil
	case <-c.ctx.Done():
		return "", c.ctx.Err()
	}
}
