package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound    = errors.New("chan not found")
	ErrUnknownTask     = errors.New("unknown task")
	ErrChannelFull     = errors.New("channel is full")
	ErrInvalidTaskFunc = errors.New("invalid task function")
)

const undefined = "undefined"

type taskItem struct {
	kind    string
	fn      interface{}
	inputs  []string
	outputs []string
}

type Conveyer struct {
	channels map[string]chan string
	size     int
	tasks    []taskItem
	mu       sync.RWMutex
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		size:     size,
		tasks:    []taskItem{},
		mu:       sync.RWMutex{},
	}
}

func (c *Conveyer) getChannel(name string) (chan string, error) {
	c.mu.RLock()

	channel, found := c.channels[name]

	c.mu.RUnlock()

	if !found {
		return nil, ErrChanNotFound
	}

	return channel, nil
}

func (c *Conveyer) createChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	channel, found := c.channels[name]

	if found {
		return channel
	}

	channel = make(chan string, c.size)
	c.channels[name] = channel

	return channel
}

func (c *Conveyer) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, channel := range c.channels {
		close(channel)
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	defer c.closeAllChannels()

	errGroup, gCtx := errgroup.WithContext(ctx)

	c.mu.RLock()

	tasksCopy := append([]taskItem(nil), c.tasks...)

	c.mu.RUnlock()

	for _, item := range tasksCopy {

		it := item

		errGroup.Go(func() error {
			return c.executeTask(gCtx, it)
		})
	}

	err := errGroup.Wait()
	if err != nil {
		return fmt.Errorf("run tasks failed: %w", err)
	}

	return nil
}

func (c *Conveyer) Send(name, data string) error {
	channel, err := c.getChannel(name)
	if err != nil {
		return ErrChanNotFound
	}

	select {
	case channel <- data:
		return nil
	default:
		return ErrChannelFull
	}
}

func (c *Conveyer) Recv(name string) (string, error) {
	channel, err := c.getChannel(name)
	if err != nil {
		return "", ErrChanNotFound
	}

	val, ok := <-channel

	if !ok {
		return undefined, nil
	}

	return val, nil
}

func (c *Conveyer) executeTask(ctx context.Context, item taskItem) error {
	switch item.kind {
	case "decorator":
		decFn, ok := item.fn.(func(context.Context, chan string, chan string) error)

		if !ok {
			return ErrInvalidTaskFunc
		}

		inputChannel := c.createChannel(item.inputs[0])

		outputChannel := c.createChannel(item.outputs[0])

		return decFn(ctx, inputChannel, outputChannel)

	case "multiplexer":
		muxFn, ok := item.fn.(func(context.Context, []chan string, chan string) error)

		if !ok {
			return ErrInvalidTaskFunc
		}

		ins := make([]chan string, len(item.inputs))

		for index, name := range item.inputs {
			ins[index] = c.createChannel(name)
		}

		outputChannel := c.createChannel(item.outputs[0])

		return muxFn(ctx, ins, outputChannel)

	case "separator":
		sepFn, ok := item.fn.(func(context.Context, chan string, []chan string) error)

		if !ok {
			return ErrInvalidTaskFunc
		}

		outs := make([]chan string, len(item.outputs))

		for index, name := range item.outputs {
			outs[index] = c.createChannel(name)
		}

		inputChannel := c.createChannel(item.inputs[0])

		return sepFn(ctx, inputChannel, outs)
	}

	return ErrUnknownTask
}

func (c *Conveyer) RegisterDecorator(
	decFunc func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	c.createChannel(input)
	c.createChannel(output)

	c.mu.Lock()
	c.tasks = append(c.tasks, taskItem{
		kind:    "decorator",
		fn:      decFunc,
		inputs:  []string{input},
		outputs: []string{output},
	})
	c.mu.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(
	decFunc func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	for _, name := range inputs {
		c.createChannel(name)
	}

	c.createChannel(output)

	c.mu.Lock()
	c.tasks = append(c.tasks, taskItem{
		kind:    "multiplexer",
		fn:      decFunc,
		inputs:  inputs,
		outputs: []string{output},
	})
	c.mu.Unlock()
}

func (c *Conveyer) RegisterSeparator(
	decFunc func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	c.createChannel(input)

	for _, name := range outputs {
		c.createChannel(name)
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, taskItem{
		kind:    "separator",
		fn:      decFunc,
		inputs:  []string{input},
		outputs: outputs,
	})
	c.mu.Unlock()
}
