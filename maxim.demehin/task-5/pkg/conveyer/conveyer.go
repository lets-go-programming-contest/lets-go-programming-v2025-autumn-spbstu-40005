package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

const undefined_msg = "undefined"

type ConveyerType struct {
	size     int
	channels map[string]chan string
	tasks    []func(ctx context.Context) error
	mutex    sync.RWMutex
}

func New(size int) *ConveyerType {
	return &ConveyerType{
		size:     size,
		channels: make(map[string]chan string),
		tasks:    make([]func(ctx context.Context) error, 0),
	}
}

func (c *ConveyerType) getOrCreateChannel(name string) chan string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *ConveyerType) RegisterDecorator(
	fn func(ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string,
	output string,
) {
	task := func(ctx context.Context) error {
		inputChan := c.getOrCreateChannel(input)

		outputChan := c.getOrCreateChannel(output)

		return fn(ctx, inputChan, outputChan)
	}

	c.tasks = append(c.tasks, task)
}

func (c *ConveyerType) RegisterMultiplexer(
	fn func(ctx context.Context,
		inputs []chan string,
		output chan string,
	) error,
	inputs []string,
	output string,
) {
	task := func(ctx context.Context) error {
		inputChans := make([]chan string, len(inputs))

		for index, inputName := range inputs {
			inputChans[index] = c.getOrCreateChannel(inputName)
		}

		outputChan := c.getOrCreateChannel(output)

		return fn(ctx, inputChans, outputChan)
	}

	c.tasks = append(c.tasks, task)
}

func (c *ConveyerType) RegisterSeparator(
	fn func(ctx context.Context,
		input chan string,
		outputs []chan string,
	) error,
	input string,
	outputs []string,
) {
	task := func(ctx context.Context) error {
		outputChans := make([]chan string, len(outputs))

		for index, outputName := range outputs {
			outputChans[index] = c.getOrCreateChannel(outputName)
		}

		inputChan := c.getOrCreateChannel(input)

		return fn(ctx, inputChan, outputChans)
	}

	c.tasks = append(c.tasks, task)
}

func (c *ConveyerType) closeChannels() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for name, ch := range c.channels {
		close(ch)

		delete(c.channels, name)
	}
}

func (c *ConveyerType) Run(ctx context.Context) error {
	defer c.closeChannels()

	group, ctx := errgroup.WithContext(ctx)

	for _, task := range c.tasks {
		group.Go(func() error {
			return task(ctx)
		})
	}

	err := group.Wait()
	if err != nil {
		return fmt.Errorf("conveyer execution failed: %w", err)
	}

	return nil
}

func (c *ConveyerType) getChannel(name string) (chan string, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	ch, exists := c.channels[name]

	return ch, exists
}

func (c *ConveyerType) Send(input string, data string) error {
	ch, exists := c.getChannel(input)

	if !exists {
		return ErrChanNotFound
	}

	ch <- data

	return nil
}

func (c *ConveyerType) Recv(output string) (string, error) {
	ch, exists := c.getChannel(output)

	if !exists {
		return "", ErrChanNotFound
	}

	data, ok := <-ch

	if !ok {
		return undefined_msg, nil
	}

	return data, nil
}
