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

type Handlers func(ctx context.Context) error

type Conveyer struct {
	mutex    sync.RWMutex
	channels map[string]chan string
	size     int
	handlers []Handlers
}

func New(size int) *Conveyer {
	return &Conveyer{
		mutex:    sync.RWMutex{},
		channels: make(map[string]chan string),
		size:     size,
		handlers: make([]Handlers, 0),
	}
}

func (c *Conveyer) getChannelOrCreate(name string) chan string {
	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch

	return ch
}

func (c *Conveyer) getChannel(name string) (chan string, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	ch, exists := c.channels[name]
	if !exists {
		return nil, fmt.Errorf("%w: channel '%s' not found", ErrChannel, name)
	}

	return ch, nil
}

func (c *Conveyer) RegisterDecorator(
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

func (c *Conveyer) RegisterMultiplexer(
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

func (c *Conveyer) RegisterSeparator(
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

func (c *Conveyer) closeAllChannels() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, ch := range c.channels {
		close(ch)
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, h := range c.handlers {
		group.Go(func() error {
			return h(ctx)
		})
	}

	err := group.Wait()

	c.closeAllChannels()

	if err != nil {
		return fmt.Errorf("conveyer execution failed: %w", err)
	}

	return nil
}

func (c *Conveyer) Send(input string, data string) error {
	ch, err := c.getChannel(input)
	if err != nil {
		return err
	}

	ch <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	ch, err := c.getChannel(output)
	if err != nil {
		return "", err
	}

	data, ok := <-ch
	if !ok {
		return undefinedValue, nil
	}

	return data, nil
}
