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
	if ch, ok := c.channels[name]; ok {
		return ch
	}
	ch := make(chan string, c.sizeChan)
	c.channels[name] = ch

	return ch
}

func (c *ConveyerStruct) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	inChan := c.getChannel(input)
	outChan := c.getChannel(output)

	task := func(ctx context.Context) error {
		return fn(ctx, inChan, outChan)
	}

	c.mute.Lock()
	c.tasks = append(c.tasks, task)
	c.mute.Unlock()
}

func (c *ConveyerStruct) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	inChans := make([]chan string, 0, len(inputs))

	for _, name := range inputs {
		inChans = append(inChans, c.getChannel(name))
	}
	outChan := c.getChannel(output)

	task := func(ctx context.Context) error {
		return fn(ctx, inChans, outChan)
	}

	c.mute.Lock()
	c.tasks = append(c.tasks, task)
	c.mute.Unlock()
}

func (c *ConveyerStruct) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	inChan := c.getChannel(input)
	outChans := make([]chan string, 0, len(outputs))

	for _, name := range outputs {
		outChans = append(outChans, c.getChannel(name))
	}

	task := func(ctx context.Context) error {
		return fn(ctx, inChan, outChans)
	}

	c.mute.Lock()
	c.tasks = append(c.tasks, task)
	c.mute.Unlock()
}

func (c *ConveyerStruct) closeChans() {
	c.mute.Lock()
	defer c.mute.Unlock()

	for _, ch := range c.channels {
		select {
		case <-ch:
		default:
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("recovered from panic: %v", r)
					}
				}()
				close(ch)
			}()
		}
	}
}

func (c *ConveyerStruct) Run(ctx context.Context) error {
	c.mute.Lock()

	group, ctx := errgroup.WithContext(ctx)

	for _, task := range c.tasks {
		group.Go(func() error {
			return task(ctx)
		})
	}

	err := group.Wait()

	c.closeChans()

	c.mute.Unlock()

	if err != nil {
		return fmt.Errorf("task run error: %w", err)
	}

	return nil
}

func (c *ConveyerStruct) Send(input string, data string) error {
	c.mute.Lock()
	ch, ok := c.channels[input]
	c.mute.Unlock()

	if !ok {
		return fmt.Errorf("error in send data: %w", ErrChanNotFound)
	}

	ch <- data

	return nil
}

func (c *ConveyerStruct) Recv(output string) (string, error) {
	c.mute.Lock()
	ch, ok := c.channels[output]
	c.mute.Unlock()

	if !ok {
		return "", fmt.Errorf("error in receiv data: %w", ErrChanNotFound)
	}

	data, ok := <-ch
	if !ok {
		return "undefined channel", nil
	}

	return data, nil
}
