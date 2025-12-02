package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

const undefinedMsg = "undefined"

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
		mutex:    sync.RWMutex{},
	}
}

func (c *ConveyerType) getOrCreateChannel(name string) chan string {
	if channel, exists := c.channels[name]; exists {
		return channel
	}

	channel := make(chan string, c.size)
	c.channels[name] = channel

	return channel
}

func (c *ConveyerType) RegisterDecorator(
	function func(ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string,
	output string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inputChan := c.getOrCreateChannel(input)
	outputChan := c.getOrCreateChannel(output)

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return function(ctx, inputChan, outputChan)
	})
}

func (c *ConveyerType) RegisterMultiplexer(
	function func(ctx context.Context,
		inputs []chan string,
		output chan string,
	) error,
	inputs []string,
	output string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	outputChan := c.getOrCreateChannel(output)
	inputChans := make([]chan string, len(inputs))

	for index, inputName := range inputs {
		inputChans[index] = c.getOrCreateChannel(inputName)
	}

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return function(ctx, inputChans, outputChan)
	})
}

func (c *ConveyerType) RegisterSeparator(
	function func(ctx context.Context,
		input chan string,
		outputs []chan string,
	) error,
	input string,
	outputs []string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inputChan := c.getOrCreateChannel(input)
	outputChans := make([]chan string, len(outputs))

	for index, outputName := range outputs {
		outputChans[index] = c.getOrCreateChannel(outputName)
	}

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return function(ctx, inputChan, outputChans)
	})
}

func (c *ConveyerType) closeChannels() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, channel := range c.channels {
		close(channel)
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
		return fmt.Errorf("error while running tasks: %w", err)
	}

	return nil
}

func (c *ConveyerType) getChannel(name string) (chan string, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	channel, exists := c.channels[name]
	if !exists {
		return nil, ErrChanNotFound
	}

	return channel, nil
}

func (c *ConveyerType) Send(input string, data string) error {
	channel, err := c.getChannel(input)
	if err != nil {
		return fmt.Errorf("error while sending data: %w", err)
	}

	channel <- data

	return nil
}

func (c *ConveyerType) Recv(output string) (string, error) {
	channel, err := c.getChannel(output)
	if err != nil {
		return "", fmt.Errorf("error whule recieving data: %w", err)
	}

	data, ok := <-channel
	if !ok {
		return undefinedMsg, nil
	}

	return data, nil
}
