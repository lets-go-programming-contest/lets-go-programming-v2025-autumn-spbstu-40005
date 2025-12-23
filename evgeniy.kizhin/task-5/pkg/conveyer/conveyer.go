package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const Undefined = "undefined"

var ErrChanNotFound = errors.New("chan not found")

type Task func(context.Context) error

type Conveyor struct {
	mu     sync.RWMutex
	chans  map[string]chan string
	tasks  []Task
	buffer int
}

func New(size int) *Conveyor {
	return &Conveyor{
		mu:     sync.RWMutex{},
		chans:  make(map[string]chan string),
		tasks:  make([]Task, 0),
		buffer: size,
	}
}

func (c *Conveyor) get(name string) (chan string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	inch, ok := c.chans[name]

	if !ok {
		return nil, ErrChanNotFound
	}

	return inch, nil
}

func (c *Conveyor) getOrCreate(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if inch, ok := c.chans[name]; ok {
		return inch
	}

	ch := make(chan string, c.buffer)
	c.chans[name] = ch

	return ch
}

func (c *Conveyor) RegisterDecorator(
	handlerFn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inch := c.getOrCreate(input)
	out := c.getOrCreate(output)
	task := func(ctx context.Context) error {
		return handlerFn(ctx, inch, out)
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *Conveyor) RegisterMultiplexer(
	handlerFn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	ins := make([]chan string, 0, len(inputs))
	for _, name := range inputs {
		ins = append(ins, c.getOrCreate(name))
	}

	out := c.getOrCreate(output)

	task := func(ctx context.Context) error {
		return handlerFn(ctx, ins, out)
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *Conveyor) RegisterSeparator(
	handlerFn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inch := c.getOrCreate(input)
	outs := make([]chan string, 0, len(outputs))

	for _, name := range outputs {
		outs = append(outs, c.getOrCreate(name))
	}

	task := func(ctx context.Context) error {
		return handlerFn(ctx, inch, outs)
	}

	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *Conveyor) Run(ctx context.Context) error {
	c.mu.RLock()
	group, gctx := errgroup.WithContext(ctx)

	for i := range c.tasks {
		task := c.tasks[i]

		group.Go(func() error {
			return task(gctx)
		})
	}
	c.mu.RUnlock()

	err := group.Wait()

	c.mu.Lock()
	for _, ch := range c.chans {
		close(ch)
	}

	c.mu.Unlock()

	if err != nil {
		return fmt.Errorf("pipeline error: %w", err)
	}

	return nil
}

func (c *Conveyor) Send(input string, data string) error {
	inch, err := c.get(input)
	if err != nil {
		return ErrChanNotFound
	}

	inch <- data

	return nil
}

func (c *Conveyor) Recv(output string) (string, error) {
	inch, err := c.get(output)
	if err != nil {
		return "", ErrChanNotFound
	}

	msg, found := <-inch
	if !found {
		return Undefined, nil
	}

	return msg, nil
}
