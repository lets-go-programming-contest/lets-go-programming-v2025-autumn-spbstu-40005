package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

type Task func(ctx context.Context) error

type ConveyerStruct struct {
	channels map[string]chan string
	mute     sync.RWMutex
	tasks    []Task
	sizeChan int
}

func New(size int) *ConveyerStruct {
	return &ConveyerStruct{
		channels: make(map[string]chan string),
		mute:     sync.RWMutex{},
		tasks:    make([]Task, 0),
		sizeChan: size,
	}
}

func (c *ConveyerStruct) getChannel(name string) (chan string, error) {
	c.mute.Lock()
	defer c.mute.Unlock()

	channel, ok := c.channels[name]
	if !ok {
		channel := make(chan string, c.sizeChan)
		c.channels[name] = channel

		return channel, ErrChanNotFound
	}

	return channel, nil
}

func (c *ConveyerStruct) RegisterDecorator(
	funct func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	inChan, _ := c.getChannel(input)
	outChan, _ := c.getChannel(output)

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
		itChan, _ := c.getChannel(name)
		inChans = append(inChans, itChan)
	}

	outChan, _ := c.getChannel(output)

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
	inChan, _ := c.getChannel(input)
	outChans := make([]chan string, 0, len(outputs))

	for _, name := range outputs {
		itChan, _ := c.getChannel(name)
		outChans = append(outChans, itChan)
	}

	task := func(ctx context.Context) error {
		return funct(ctx, inChan, outChans)
	}

	c.mute.Lock()
	c.tasks = append(c.tasks, task)
	c.mute.Unlock()
}

func (c *ConveyerStruct) closeChans() {
	for _, channel := range c.channels {
		func() {
			defer func() {
				_ = recover()
			}()
			close(channel)
		}()
	}
}

func (c *ConveyerStruct) Run(ctx context.Context) error {
	group, gCtx := errgroup.WithContext(ctx)

	c.mute.RLock()
	tasks := c.tasks
	c.mute.RUnlock()

	for _, task := range tasks {
		t := task

		group.Go(func() error {
			return t(gCtx)
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
	channel, availability := c.channels[input]
	c.mute.RUnlock()

	if !availability {
		return fmt.Errorf("error in send data: %w", ErrChanNotFound)
	}

	defer func() {
		_ = recover()
	}()
	channel <- data

	return nil
}

func (c *ConveyerStruct) Recv(output string) (string, error) {
	c.mute.RLock()
	channel, availability := c.channels[output]
	c.mute.RUnlock()

	if !availability {
		return "", fmt.Errorf("error in receiv data: %w", ErrChanNotFound)
	}

	data, ok := <-channel
	if !ok {
		return "undefined", nil
	}

	return data, nil
}
