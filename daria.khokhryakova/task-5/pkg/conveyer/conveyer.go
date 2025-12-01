package conveyer

import (
	"context"
	"errors"
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
		channels: make(map[string]chan string),
		size:     size,
		handlers: make([]func(ctx context.Context) error, 0),
	}
}

func (c *conveyerImpl) getChannelOrCreate(name string) chan string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

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
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inputChan := c.getChannelOrCreate(input)
	outputChan := c.getChannelOrCreate(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputChan, outputChan)
	})
}

func (c *conveyerImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	outputChan := c.getChannelOrCreate(output)
	inputChans := make([]chan string, len(inputs))
	for i, input := range inputs {
		inputChans[i] = c.getChannelOrCreate(input)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputChans, outputChan)
	})
}

func (c *conveyerImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inputChan := c.getChannelOrCreate(input)
	outputChans := make([]chan string, len(outputs))
	for i, output := range outputs {
		outputChans[i] = c.getChannelOrCreate(output)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputChan, outputChans)
	})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, handler := range c.handlers {
		h := handler
		g.Go(func() error {
			return h(ctx)
		})
	}

	err := g.Wait()

	c.mutex.Lock()
	for _, ch := range c.channels {
		close(ch)
	}
	c.channels = make(map[string]chan string)
	c.handlers = make([]func(ctx context.Context) error, 0)
	c.mutex.Unlock()

	return err
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
