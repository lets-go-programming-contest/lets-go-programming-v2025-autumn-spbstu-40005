package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound    = errors.New("chan not found")
	ErrConveyerRunning = errors.New("conveyer is already running")
	ErrConveyerStopped = errors.New("conveyer is not running")
	ErrChannelFull     = errors.New("channel is full")
	ErrNoDataAvailable = errors.New("no data available")
	ErrUndefined       = "undefined"
)

type Conveyer struct {
	mu       sync.RWMutex
	channels map[string]chan string
	size     int
	tasks    []task
	running  bool
	cancel   context.CancelFunc
}

type task struct {
	fn      func(ctx context.Context) error
	inputs  []string
	outputs []string
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		size:     size,
	}
}

func (c *Conveyer) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inputChan := c.getOrCreateChannel(input)
	outputChan := c.getOrCreateChannel(output)

	taskFn := func(ctx context.Context) error {
		return fn(ctx, inputChan, outputChan)
	}

	c.tasks = append(c.tasks, task{
		fn:      taskFn,
		inputs:  []string{input},
		outputs: []string{output},
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputChans := make([]chan string, len(inputs))
	for i, input := range inputs {
		inputChans[i] = c.getOrCreateChannel(input)
	}
	outputChan := c.getOrCreateChannel(output)

	taskFn := func(ctx context.Context) error {
		return fn(ctx, inputChans, outputChan)
	}

	c.tasks = append(c.tasks, task{
		fn:      taskFn,
		inputs:  inputs,
		outputs: []string{output},
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inputChan := c.getOrCreateChannel(input)
	outputChans := make([]chan string, len(outputs))
	for i, output := range outputs {
		outputChans[i] = c.getOrCreateChannel(output)
	}

	taskFn := func(ctx context.Context) error {
		return fn(ctx, inputChan, outputChans)
	}

	c.tasks = append(c.tasks, task{
		fn:      taskFn,
		inputs:  []string{input},
		outputs: outputs,
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return ErrConveyerRunning
	}
	c.running = true

	runCtx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.running = false
		c.cancel = nil
		c.closeAllChannels()
		c.mu.Unlock()
	}()

	g, gCtx := errgroup.WithContext(runCtx)

	for _, task := range c.tasks {
		t := task
		g.Go(func() error {
			return t.fn(gCtx)
		})
	}

	if err := g.Wait(); err != nil {
		cancel()
		return err
	}

	return nil
}

func (c *Conveyer) closeAllChannels() {
	for name, ch := range c.channels {
		close(ch)
		delete(c.channels, name)
	}
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, exists := c.channels[input]
	if !exists {
		return ErrChanNotFound
	}

	if !c.running {
		return ErrConveyerStopped
	}

	select {
	case ch <- data:
		return nil
	default:
		return ErrChannelFull
	}
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, exists := c.channels[output]
	running := c.running
	c.mu.RUnlock()

	if !exists {
		return "", ErrChanNotFound
	}

	if !running {
		return "", ErrConveyerStopped
	}

	select {
	case data, ok := <-ch:
		if !ok {
			return ErrUndefined, nil
		}
		return data, nil
	default:
		return "", ErrNoDataAvailable
	}
}

func (c *Conveyer) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cancel != nil {
		c.cancel()
	}
	c.closeAllChannels()
	c.running = false
}
