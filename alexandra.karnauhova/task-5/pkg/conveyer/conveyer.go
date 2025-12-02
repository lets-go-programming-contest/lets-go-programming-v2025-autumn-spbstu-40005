package conveyer

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
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
	in := c.getChannel(input)
	out := c.getChannel(output)

	task := func(ctx context.Context) error {
		return fn(ctx, in, out)
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
	var inChans []chan string
	for _, name := range inputs {
		inChans = append(inChans, c.getChannel(name))
	}
	out := c.getChannel(output)

	task := func(ctx context.Context) error {
		return fn(ctx, inChans, out)
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
	in := c.getChannel(input)
	var outChans []chan string
	for _, name := range outputs {
		outChans = append(outChans, c.getChannel(name))
	}

	task := func(ctx context.Context) error {
		return fn(ctx, in, outChans)
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
				defer func() { recover() }()
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
		return fmt.Errorf("error while running tasks: %w", err)
	}

	return nil
}
