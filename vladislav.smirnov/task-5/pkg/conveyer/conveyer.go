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

type task func(context.Context) error

type Conveyer struct {
	channels map[string]chan string
	size     int
	tasks    []task
	mu       sync.RWMutex
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		size:     size,
		tasks:    []task{},
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

	for _, task := range c.tasks {
		currentTask := task

		errGroup.Go(func() error {
			return currentTask(gCtx)
		})
	}

	c.mu.RUnlock()

	err := errGroup.Wait()

	if err != nil {
		return fmt.Errorf("run tasks failed: %w", err)
	}

	return nil
}

func (c *Conveyer) Send(name, data string) error {
	channel, err := c.getChannel(name)
	if err != nil {
		return err
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

func (c *Conveyer) RegisterDecorator(
	decFunc func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	inChannel := c.createChannel(input)
	outChannel := c.createChannel(output)

	c.mu.Lock()
	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return decFunc(ctx, inChannel, outChannel)
	})
	c.mu.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(
	muxFunc func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	ins := make([]chan string, len(inputs))
	for i, name := range inputs {
		ins[i] = c.createChannel(name)
	}

	outChannel := c.createChannel(output)

	c.mu.Lock()
	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return muxFunc(ctx, ins, outChannel)
	})
	c.mu.Unlock()
}

func (c *Conveyer) RegisterSeparator(
	sepFunc func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	inChannel := c.createChannel(input)
	outs := make([]chan string, len(outputs))

	for i, name := range outputs {
		outs[i] = c.createChannel(name)
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return sepFunc(ctx, inChannel, outs)
	})
	c.mu.Unlock()
}
