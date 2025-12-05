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
		fn func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)

	RegisterMultiplexer(
		fn func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)

	RegisterSeparator(
		fn func(ctx context.Context, input chan string, outputs []chan string) error,
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

func New(size int) Conveyer {
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

	ch, exists := c.channels[name]
	return ch, exists
}

func (c *conveyerImpl) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.channelBuffer)
	c.channels[name] = ch
	return ch
}

func (c *conveyerImpl) Send(input string, data string) error {
	ch, exists := c.getChannel(input)
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
	case ch <- data:
		return nil
	default:
		return ErrSendFailed
	}
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	ch, exists := c.getChannel(output)
	if !exists {
		return "", fmt.Errorf("receive failed: %w", ErrChanNotFound)
	}

	value, ok := <-ch
	if !ok {
		return undefinedValue, nil
	}

	return value, nil
}

func (c *conveyerImpl) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	srcChan := c.getOrCreateChannel(input)
	dstChan := c.getOrCreateChannel(output)

	task := func(ctx context.Context) error {
		return fn(ctx, srcChan, dstChan)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.tasks = append(c.tasks, task)
}

func (c *conveyerImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	srcChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		srcChans[i] = c.getOrCreateChannel(name)
	}

	dstChan := c.getOrCreateChannel(output)

	task := func(ctx context.Context) error {
		return fn(ctx, srcChans, dstChan)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.tasks = append(c.tasks, task)
}

func (c *conveyerImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	srcChan := c.getOrCreateChannel(input)

	dstChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		dstChans[i] = c.getOrCreateChannel(name)
	}

	task := func(ctx context.Context) error {
		return fn(ctx, srcChan, dstChans)
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
	for _, ch := range c.channels {
		if ch != nil {
			close(ch)
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
