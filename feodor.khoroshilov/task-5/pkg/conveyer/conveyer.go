package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const undefinedValue = "undefined"

var (
	ErrChannelNotFound    = errors.New("chan not found")
	ErrContextCannotBeNil = errors.New("context cannot be nil")
)

type WorkerFunc func(ctx context.Context) error

type conveyer struct {
	bufferSize int
	channels   map[string]chan string
	workers    []WorkerFunc
	mutex      sync.RWMutex
}

func New(size int) *conveyer {
	return &conveyer{
		bufferSize: size,
		channels:   make(map[string]chan string),
		workers:    make([]WorkerFunc, 0),
		mutex:      sync.RWMutex{},
	}
}

func (c *conveyer) getChannelOrCreateLocked(name string) chan string {
	if ch, ok := c.channels[name]; ok {
		return ch
	}

	ch := make(chan string, c.bufferSize)
	c.channels[name] = ch

	return ch
}

func (c *conveyer) getChannel(name string) (chan string, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	ch, ok := c.channels[name]
	if !ok {
		return nil , ErrChannelNotFound
	}

	return ch, nil
}

func (c *conveyer) RegisterDecorator(
	processor func(ctx context.Context, in chan string, out chan string) error,
	inName,
	outName string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inChan := c.getChannelOrCreateLocked(inName)
	outChan := c.getChannelOrCreateLocked(outName)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return processor(ctx, inChan, outChan)
	})
}

func (c *conveyer) RegisterMultiplexer(
	processor func(ctx context.Context, ins []chan string, out chan string) error,
	inNames []string,
	outName string,
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

func (c *conveyer) RegisterSeparator(
	processor func(ctx context.Context, in chan string, outs []chan string) error,
	inName string,
	outNames []string,
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

func (c *conveyer) Run(ctx context.Context) error {
	if ctx == nil {
		return ErrContextCannotBeNil
	}

	c.mutex.RLock()
	workers := make([]WorkerFunc, len(c.workers))
	copy(workers, c.workers)
	c.mutex.RUnlock()

	errGroup, ctx := errgroup.WithContext(ctx)

	for _, w := range workers {
		errGroup.Go(func() error {
			return w(ctx)
		})
	}

	err := errGroup.Wait()

	c.closeAll()

	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

func (c *conveyer) closeAll() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, ch := range c.channels {
		close(ch)
	}
}

func (c *conveyer) Send(name string, data string) error {
	ch, exists := c.getChannel(name)
	if !exists {
		return ErrChannelNotFound
	}

	ch <- data

	return nil
}

func (c *conveyer) Recv(name string) (string, error) {
	ch, exists := c.getChannel(name)
	if !exists {
		return "", ErrChannelNotFound
	}

	value, ok := <-ch
	if !ok {
		return undefinedValue, nil
	}

	return value, nil
}
