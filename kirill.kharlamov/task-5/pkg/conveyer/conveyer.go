package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const (
	emptyValue = "undefined"
)

var ErrChannelNotFound = errors.New("channel not found")

type Conveyer struct {
	bufferSize int
	channels   map[string]chan string
	tasks      []func(context.Context) error
	mutex      sync.RWMutex
}

func New(size int) *Conveyer {
	return &Conveyer{
		bufferSize: size,
		channels:   make(map[string]chan string),
		tasks:      make([]func(context.Context) error, 0),
		mutex:      sync.RWMutex{},
	}
}

func (c *Conveyer) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	inName, outName string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inCh := c.ensureChannel(inName)
	outCh := c.ensureChannel(outName)

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inNames []string,
	outName string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	outCh := c.ensureChannel(outName)
	inChs := make([]chan string, len(inNames))

	for i, name := range inNames {
		inChs[i] = c.ensureChannel(name)
	}

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return fn(ctx, inChs, outCh)
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	inName string,
	outNames []string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inCh := c.ensureChannel(inName)
	outChs := make([]chan string, len(outNames))

	for i, name := range outNames {
		outChs[i] = c.ensureChannel(name)
	}

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return fn(ctx, inCh, outChs)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	errGroup, ctx := errgroup.WithContext(ctx)

	for _, task := range c.tasks {
		task := task
		errGroup.Go(func() error {
			return task(ctx)
		})
	}

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("conveyer execution failed: %w", err)
	}

	c.closeAllChannels()
	return nil
}

func (c *Conveyer) closeAllChannels() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, ch := range c.channels {
		close(ch)
	}
}

func (c *Conveyer) Send(channelName string, data string) error {
	c.mutex.RLock()
	ch, exists := c.channels[channelName]
	c.mutex.RUnlock()

	if !exists {
		return ErrChannelNotFound
	}

	select {
	case ch <- data:
		return nil
	case <-context.TODO().Done():
		return context.Canceled
	}
}

func (c *Conveyer) Recv(channelName string) (string, error) {
	c.mutex.RLock()
	ch, exists := c.channels[channelName]
	c.mutex.RUnlock()

	if !exists {
		return "", ErrChannelNotFound
	}

	select {
	case val, ok := <-ch:
		if !ok {
			return emptyValue, nil
		}
		return val, nil
	case <-context.TODO().Done():
		return "", context.Canceled
	}
}

func (c *Conveyer) ensureChannel(name string) chan string {
	c.mutex.RLock()
	if ch, exists := c.channels[name]; exists {
		c.mutex.RUnlock()
		return ch
	}
	c.mutex.RUnlock()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.bufferSize)
	c.channels[name] = ch
	return ch
}
