package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const undefinedValue = "undefined"

var (
	ErrChanNotFound = errors.New("chan not found")
	ErrSendFailed   = errors.New("send failed: channel closed or full")
)

type Conveyer interface {
	RegisterDecorator(
		decoratorFunc func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)

	RegisterMultiplexer(
		multiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)

	RegisterSeparator(
		separatorFunc func(ctx context.Context, input chan string, outputs []chan string) error,
		input string,
		outputs []string,
	)

	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type conveyerImpl struct {
	mu            sync.RWMutex
	channels      map[string]chan string
	tasks         []func(context.Context) error
	channelBuffer int
	closed        bool
}

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		mu:            sync.RWMutex{},
		channels:      make(map[string]chan string),
		tasks:         make([]func(context.Context) error, 0),
		channelBuffer: size,
		closed:        false,
	}
}

func (c *conveyerImpl) getChannel(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	channel, exists := c.channels[name]

	return channel, exists
}

func (c *conveyerImpl) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if channel, exists := c.channels[name]; exists {
		return channel
	}

	channel := make(chan string, c.channelBuffer)
	c.channels[name] = channel

	return channel
}

func (c *conveyerImpl) Send(input string, data string) error {
	channel, exists := c.getChannel(input)
	if !exists {
		return fmt.Errorf("send failed: %w", ErrChanNotFound)
	}

	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()

		return ErrSendFailed
	}
	c.mu.RUnlock()

	select {
	case channel <- data:
		return nil
	default:
		return ErrSendFailed
	}
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	channel, exists := c.getChannel(output)
	if !exists {
		return "", fmt.Errorf("receive failed: %w", ErrChanNotFound)
	}

	value, ok := <-channel
	if !ok {
		return undefinedValue, nil
	}

	return value, nil
}

func (c *conveyerImpl) RegisterDecorator(
	decoratorFunc func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	srcChan := c.getOrCreateChannel(input)
	dstChan := c.getOrCreateChannel(output)

	task := func(ctx context.Context) error {
		return decoratorFunc(ctx, srcChan, dstChan)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.tasks = append(c.tasks, task)
}

func (c *conveyerImpl) RegisterMultiplexer(
	multiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	srcChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		srcChans[i] = c.getOrCreateChannel(name)
	}

	dstChan := c.getOrCreateChannel(output)

	task := func(ctx context.Context) error {
		return multiplexerFunc(ctx, srcChans, dstChan)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.tasks = append(c.tasks, task)
}

func (c *conveyerImpl) RegisterSeparator(
	separatorFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	srcChan := c.getOrCreateChannel(input)

	dstChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		dstChans[i] = c.getOrCreateChannel(name)
	}

	task := func(ctx context.Context) error {
		return separatorFunc(ctx, srcChan, dstChans)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.tasks = append(c.tasks, task)
}

func (c *conveyerImpl) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}

	c.closed = true
	for _, channel := range c.channels {
		if channel != nil {
			close(channel)
		}
	}
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	defer c.closeAllChannels()

	group, groupCtx := errgroup.WithContext(ctx)

	c.mu.RLock()
	tasks := make([]func(context.Context) error, len(c.tasks))
	copy(tasks, c.tasks)
	c.mu.RUnlock()

	for _, task := range tasks {
		currentTask := task

		group.Go(func() error {
			return currentTask(groupCtx)
		})
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("conveyer execution terminated: %w", err)
	}

	return nil
}
