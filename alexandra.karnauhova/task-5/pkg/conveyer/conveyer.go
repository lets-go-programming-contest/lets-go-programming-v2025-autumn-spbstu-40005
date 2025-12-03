package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound = errors.New("chan not found")
)

type ConveyerStruct struct {
	channels map[string]chan string
	mute     sync.RWMutex
	tasks    []func(context.Context) error
	sizeChan int
}

func New(size int) *ConveyerStruct {
	return &ConveyerStruct{
		channels: make(map[string]chan string),
		mute:     sync.RWMutex{},
		tasks:    make([]func(ctx context.Context) error, 0),
		sizeChan: size,
	}
}

func (c *ConveyerStruct) getChannel(name string) chan string {
	c.mute.Lock()
	defer c.mute.Unlock()

	if channel, ok := c.channels[name]; ok {
		return channel
	}

	channel := make(chan string, c.sizeChan)
	c.channels[name] = channel

	return channel
}

func (c *ConveyerStruct) RegisterDecorator(
	funct func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	inChan := c.getChannel(input)
	outChan := c.getChannel(output)

	task := func(ctx context.Context) error {
		return funct(ctx, inChan, outChan)
	}

	c.mute.Lock()
	c.tasks = append(c.tasks, task)
	c.mute.Unlock()
}

func (c *ConveyerStruct) RegisterMultiplexer(
	funct func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	inChans := make([]chan string, 0, len(inputs))

	for _, name := range inputs {
		inChans = append(inChans, c.getChannel(name))
	}

	outChan := c.getChannel(output)

	task := func(ctx context.Context) error {
		return funct(ctx, inChans, outChan)
	}

	c.mute.Lock()
	c.tasks = append(c.tasks, task)
	c.mute.Unlock()
}

func (c *ConveyerStruct) RegisterSeparator(
	funct func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	inChan := c.getChannel(input)
	outChans := make([]chan string, 0, len(outputs))

	for _, name := range outputs {
		outChans = append(outChans, c.getChannel(name))
	}

	task := func(ctx context.Context) error {
		return funct(ctx, inChan, outChans)
	}

	c.mute.Lock()
	c.tasks = append(c.tasks, task)
	c.mute.Unlock()
}

func (c *ConveyerStruct) closeChans() {
	c.mute.Lock()
	defer c.mute.Unlock()

	for _, channel := range c.channels {
		close(channel)
	}
}

func (c *ConveyerStruct) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, task := range c.tasks {
		group.Go(func() error {
			return task(ctx)
		})
	}

	err := group.Wait()

	c.mute.Lock()
	c.closeChans()
	c.mute.Unlock()

	if err != nil {
		return fmt.Errorf("task run error: %w", err)
	}

	return nil
}

func (c *ConveyerStruct) Send(input string, data string) error {
	c.mute.RLock()
	channel, ok := c.channels[input]
	c.mute.RUnlock()

	if !ok {
		return fmt.Errorf("error in send data: %w", ErrChanNotFound)
	}

	channel <- data

	return nil
}

func (c *ConveyerStruct) Recv(output string) (string, error) {
	c.mute.RLock()
	channel, ok := c.channels[output]
	c.mute.RUnlock()

	if !ok {
		return "", fmt.Errorf("error in receiv data: %w", ErrChanNotFound)
	}

	data, ok := <-channel
	if !ok {
		return "undefined channel", nil
	}

	return data, nil
}
