package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const undefinedValue = "undefined"

var ErrChannelNotFound = errors.New("chan not found")

type conveyer struct {
	bufferSize int
	channels   map[string]chan string
	workers    []func(ctx context.Context) error
	mutex      sync.RWMutex
}

func New(size int) *conveyer {
	return &conveyer{
		bufferSize: size,
		channels:   make(map[string]chan string),
		workers:    make([]func(ctx context.Context) error, 0),
		mutex:      sync.RWMutex{},
	}
}

func (c *conveyer) getChannel(name string) (chan string, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	ch, ok := c.channels[name]
	return ch, ok
}

func (c *conveyer) getChannelOrCreate(name string) chan string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	if ch, ok := c.channels[name]; ok {
		return ch
	}
	
	ch := make(chan string, c.bufferSize)
	c.channels[name] = ch
	return ch
}

func (c *conveyer) RegisterDecorator(
	processor func(ctx context.Context, in chan string, out chan string) error,
	inName,
	outName string,
) {
	inChan := c.getChannelOrCreate(inName)
	outChan := c.getChannelOrCreate(outName)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.workers = append(c.workers, func(ctx context.Context) error {
		return processor(ctx, inChan, outChan)
	})
}

func (c *conveyer) RegisterMultiplexer(
	processor func(ctx context.Context, ins []chan string, out chan string) error,
	inNames []string,
	outName string,
) {
	inChans := make([]chan string, len(inNames))
	for i, name := range inNames {
		inChans[i] = c.getChannelOrCreate(name)
	}
	outChan := c.getChannelOrCreate(outName)

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.workers = append(c.workers, func(ctx context.Context) error {
		return processor(ctx, inChans, outChan)
	})
}

func (c *conveyer) RegisterSeparator(
	processor func(ctx context.Context, in chan string, outs []chan string) error,
	inName string,
	outNames []string,
) {
	inChan := c.getChannelOrCreate(inName)
	outChans := make([]chan string, len(outNames))
	for i, name := range outNames {
		outChans[i] = c.getChannelOrCreate(name)
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.workers = append(c.workers, func(ctx context.Context) error {
		return processor(ctx, inChan, outChans)
	})
}

func (c *conveyer) Run(ctx context.Context) error {
	if ctx == nil {
		return errors.New("context cannot be nil")
	}
	
	defer c.closeAll()

	errGroup, ctx := errgroup.WithContext(ctx)

	for _, w := range c.workers {
		w := w
		errGroup.Go(func() error {
			return w(ctx)
		})
	}

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

func (c *conveyer) closeAll() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for name, ch := range c.channels {
		close(ch)
		delete(c.channels, name)
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