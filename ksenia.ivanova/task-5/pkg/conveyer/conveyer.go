package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChannelNotFound = errors.New("chan not found")

type WorkerFunc func(ctx context.Context) error

type Conveyer struct {
	size     int
	channels map[string]chan string
	workers  []WorkerFunc
	mutex    sync.RWMutex
}

func (c *Conveyer) getChannelOrCreateLocked(name string) chan string {
	if ch, ok := c.channels[name]; ok {
		return ch
	}

	ch := make(chan string, c.size)
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

func (c *Conveyer) closeChannels() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, ch := range c.channels {
		close(ch)
	}
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		workers:  make([]WorkerFunc, 0),
		mutex:    sync.RWMutex{},
	}
}

func (c *Conveyer) RegisterDecorator(
	processor func(ctx context.Context, in chan string, out chan string) error,
	inName string, outName string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inChan := c.getChannelOrCreateLocked(inName)
	outChan := c.getChannelOrCreateLocked(outName)

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
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return processor(ctx, inChan, outChans)
	})
}

func (c *Conveyer) Run(parentCtx context.Context) error {
	defer c.closeChannels()

	errGroup, ctx := errgroup.WithContext(parentCtx)

	c.mutex.RLock()

	for _, work := range c.workers {
		w := work

		errGroup.Go(func() error { return w(ctx) })
	}

	c.mutex.RUnlock()

	err := errGroup.Wait()
	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

func (c *Conveyer) Send(name string, data string) error {
	ch, err := c.getChannel(name)
	if err != nil {
		return err
	}

	ch <- data

	return nil
}

func (c *Conveyer) Recv(name string) (string, error) {
	ch, err := c.getChannel(name)
	if err != nil {
		return "", err
	}

	data, ok := <-ch
	if !ok {
		return "no data", nil
	}

	return data, nil
}
