package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrChannelNotFound = errors.New("channel not found")
	ErrConveyerRunning = errors.New("conveyer is already running")
)

type WorkerFunc func(ctx context.Context) error

type Conveyer struct {
	bufferSize int
	channels   map[string]chan string
	workers    []WorkerFunc
	outputs    map[string]struct{}
	mutex      sync.RWMutex
	running    bool
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	errChan    chan error
}

func New(size int) *Conveyer {
	return &Conveyer{
		bufferSize: size,
		channels:   make(map[string]chan string),
		workers:    make([]WorkerFunc, 0),
		outputs:    make(map[string]struct{}),
		mutex:      sync.RWMutex{},
		errChan:    make(chan error, 1),
	}
}

func (c *Conveyer) getChannelOrCreateLocked(name string) chan string {
	if ch, ok := c.channels[name]; ok {
		return ch
	}

	ch := make(chan string, c.bufferSize)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) getChannel(name string) (chan string, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	ch, ok := c.channels[name]
	if !ok {
		return nil, ErrChannelNotFound
	}

	return ch, nil
}

func (c *Conveyer) RegisterDecorator(
	processor func(ctx context.Context, in chan string, out chan string) error,
	inName, outName string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inChan := c.getChannelOrCreateLocked(inName)
	outChan := c.getChannelOrCreateLocked(outName)
	c.outputs[outName] = struct{}{}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return processor(ctx, inChan, outChan)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	processor func(ctx context.Context, ins []chan string, out chan string) error,
	inNames []string, outName string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inChans := make([]chan string, len(inNames))
	for i, name := range inNames {
		inChans[i] = c.getChannelOrCreateLocked(name)
	}

	outChan := c.getChannelOrCreateLocked(outName)
	c.outputs[outName] = struct{}{}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return processor(ctx, inChans, outChan)
	})
}

func (c *Conveyer) RegisterSeparator(
	processor func(ctx context.Context, in chan string, outs []chan string) error,
	inName string, outNames []string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inChan := c.getChannelOrCreateLocked(inName)

	outChans := make([]chan string, len(outNames))
	for i, name := range outNames {
		outChans[i] = c.getChannelOrCreateLocked(name)
		c.outputs[name] = struct{}{}
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return processor(ctx, inChan, outChans)
	})
}

func (c *Conveyer) Run(parentCtx context.Context) error {
	c.mutex.Lock()
	if c.running {
		c.mutex.Unlock()
		return ErrConveyerRunning
	}
	c.running = true
	c.ctx, c.cancel = context.WithCancel(parentCtx)

	for name, ch := range c.channels {
		if _, ok := c.outputs[name]; !ok {
			close(ch)
		}
	}

	workers := make([]WorkerFunc, len(c.workers))
	copy(workers, c.workers)
	c.mutex.Unlock()

	defer func() {
		c.mutex.Lock()
		c.running = false
		c.mutex.Unlock()
	}()

	for _, w := range workers {
		c.wg.Add(1)
		go func(worker WorkerFunc) {
			defer c.wg.Done()
			if err := worker(c.ctx); err != nil {
				select {
				case c.errChan <- err:
				default:
				}
				if c.cancel != nil {
					c.cancel()
				}
			}
		}(w)
	}

	c.wg.Wait()

	c.mutex.Lock()
	for name, ch := range c.channels {
		if _, ok := c.outputs[name]; ok {
			close(ch)
		}
	}
	c.mutex.Unlock()

	select {
	case err := <-c.errChan:
		return fmt.Errorf("execution failed: %w", err)
	default:
		return nil
	}
}

func (c *Conveyer) Send(name string, data string) error {
	ch, err := c.getChannel(name)
	if err != nil {
		return err
	}

	c.mutex.RLock()
	started := c.running
	ctx := c.ctx
	c.mutex.RUnlock()

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

func (c *Conveyer) Recv(name string) (string, error) {
	ch, err := c.getChannel(name)
	if err != nil {
		return "", err
	}

	c.mutex.RLock()
	started := c.running
	ctx := c.ctx
	c.mutex.RUnlock()

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
