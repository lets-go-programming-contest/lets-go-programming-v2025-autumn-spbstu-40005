package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChannel = errors.New("chan not found")

const undefinedValue = "undefined"

type conveyerImpl struct {
	mutex    sync.RWMutex
	channels map[string]chan string
	size     int
	handlers []func(ctx context.Context) error
}

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		mutex:    sync.RWMutex{},
		channels: make(map[string]chan string),
		size:     size,
		handlers: make([]func(ctx context.Context) error, 0),
	}
}

func (c *conveyerImpl) getChannelOrCreate(name string) chan string {
	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch

	return ch
}

func (c *conveyerImpl) getChannel(name string) (chan string, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	ch, exists := c.channels[name]

	return ch, exists
}

func (c *conveyerImpl) RegisterDecorator(
	function func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inputChan := c.getChannelOrCreate(input)
	outputChan := c.getChannelOrCreate(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return function(ctx, inputChan, outputChan)
	})
}

func (c *conveyerImpl) RegisterMultiplexer(
	function func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	outputChan := c.getChannelOrCreate(output)
	inputChans := make([]chan string, len(inputs))

	for i, input := range inputs {
		inputChans[i] = c.getChannelOrCreate(input)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return function(ctx, inputChans, outputChan)
	})
}

func (c *conveyerImpl) RegisterSeparator(
	function func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inputChan := c.getChannelOrCreate(input)
	outputChans := make([]chan string, len(outputs))

	for i, output := range outputs {
		outputChans[i] = c.getChannelOrCreate(output)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return function(ctx, inputChan, outputChans)
	})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, h := range c.handlers {
		group.Go(func() error {
			return h(ctx)
		})
	}

	err := group.Wait()

	c.mutex.Lock()
	for _, ch := range c.channels {
		close(ch)
	}
	c.mutex.Unlock()

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (c *conveyerImpl) Send(input string, data string) error {
	ch, exists := c.getChannel(input)
	if !exists {
		return ErrChannel
	}

	ch <- data

	return nil
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	ch, exists := c.getChannel(output)
	if !exists {
		return "", ErrChannel
	}

	data, ok := <-ch
	if !ok {
		return undefinedValue, nil
	}

	return data, nil
}
