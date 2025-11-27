package conveyer

import (
	"context"
	"sync"
)

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
