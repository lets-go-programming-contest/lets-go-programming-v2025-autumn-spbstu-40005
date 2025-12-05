package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

const undefined = "undefined"

var ErrChannelNotFound = errors.New("chan not found")

type conveyer struct {
	mu         sync.RWMutex
	bufferSize int
	channels   map[string]chan string
	workers    []func(ctx context.Context) error
}

func New(size int) *conveyer {
	return &conveyer{
		bufferSize: size,
		channels:   make(map[string]chan string),
		workers:    make([]func(ctx context.Context) error, 0),
	}
}

func (c *conveyer) getChannel(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	ch, ok := c.channels[name]
	return ch, ok
}

func (c *conveyer) ensureChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.channels[name]; ok {
		return ch
	}

	ch := make(chan string, c.bufferSize)
	c.channels[name] = ch
	return ch
}

func (c *conveyer) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	inputName string,
	outputName string,
) {
	inputChan := c.ensureChannel(inputName)
	outputChan := c.ensureChannel(outputName)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, inputChan, outputChan)
	})
}

func (c *conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputNames []string,
	outputName string,
) {
	inputChans := make([]chan string, len(inputNames))
	for i, name := range inputNames {
		inputChans[i] = c.ensureChannel(name)
	}

	outputChan := c.ensureChannel(outputName)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, inputChans, outputChan)
	})
}

func (c *conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	inputName string,
	outputNames []string,
) {
	inputChan := c.ensureChannel(inputName)

	outputChans := make([]chan string, len(outputNames))
	for i, name := range outputNames {
		outputChans[i] = c.ensureChannel(name)
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, inputChan, outputChans)
	})
}

func (c *conveyer) Run(ctx context.Context) error {
	defer c.closeAllChannels()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	for _, worker := range c.workers {
		w := worker 
		g.Go(func() error {
			return w(ctx)
		})
	}

	return g.Wait()
}

func (c *conveyer) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for name, ch := range c.channels {
		close(ch)
		delete(c.channels, name)
	}
}

func (c *conveyer) Send(inputName string, data string) error {
	ch, exists := c.getChannel(inputName)
	if !exists {
		return ErrChannelNotFound
	}

	ch <- data
	return nil
}

func (c *conveyer) Recv(outputName string) (string, error) {
	ch, exists := c.getChannel(outputName)
	if !exists {
		return "", ErrChannelNotFound
	}

	data, ok := <-ch
	if !ok {
		return undefined, nil
	}

	return data, nil
}