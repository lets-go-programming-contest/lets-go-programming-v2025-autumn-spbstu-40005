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
	input string,
	output string,
) {
	inChan := c.ensureChannel(input)
	outChan := c.ensureChannel(output)
	
	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, inChan, outChan)
	})
}

func (c *conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputChans := make([]chan string, len(inputs))
	for i, inputName := range inputs {
		inputChans[i] = c.ensureChannel(inputName)
	}
	outChan := c.ensureChannel(output)
	
	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, inputChans, outChan)
	})
}

func (c *conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inChan := c.ensureChannel(input)
	outputChans := make([]chan string, len(outputs))
	for i, outputName := range outputs {
		outputChans[i] = c.ensureChannel(outputName)
	}
	
	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, inChan, outputChans)
	})
}

func (c *conveyer) Run(ctx context.Context) error {
	defer c.closeAllChannels()
	
	if ctx == nil {
		return errors.New("context cannot be nil")
	}
	
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
	
	for _, ch := range c.channels {
		close(ch)
	}
	c.channels = make(map[string]chan string)
}

func (c *conveyer) Send(input string, data string) error {
	ch, exists := c.getChannel(input)
	if !exists {
		return ErrChannelNotFound
	}
	
	ch <- data
	return nil
}

func (c *conveyer) Recv(output string) (string, error) {
	ch, exists := c.getChannel(output)
	if !exists {
		return "", ErrChannelNotFound
	}
	
	data, ok := <-ch
	if !ok {
		return undefined, nil
	}
	
	return data, nil
}