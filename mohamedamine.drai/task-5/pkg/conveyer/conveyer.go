package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChannelMissing = errors.New("we cant found chan")

const emptyValue = "undefined"

type Conveyor interface {
	AddDecorator(
		handler func(context.Context, chan string, chan string) error,
		in string,
		out string,
	)
	AddMultiplexer(
		handler func(context.Context, []chan string, chan string) error,
		inputs []string,
		output string,
	)
	AddSeparator(
		handler func(context.Context, chan string, []chan string) error,
		input string,
		outputs []string,
	)

	Run(context.Context) error
	Push(string, string) error
	Pull(string) (string, error)
}

type Flow struct {
	bufSize  int
	bus      map[string]chan string
	tasks    []func(context.Context) error
	busGuard sync.RWMutex
}

func New(size int) *Flow {
	return &Flow{
		bufSize:  size,
		bus:      make(map[string]chan string),
		tasks:    make([]func(context.Context) error, 0),
		busGuard: sync.RWMutex{},
	}
}

func (f *Flow) ensure(name string) chan string {
	f.busGuard.Lock()
	defer f.busGuard.Unlock()

	if c, ok := f.bus[name]; ok {
		return c
	}

	ch := make(chan string, f.bufSize)
	f.bus[name] = ch
	return ch
}

func (f *Flow) AddDecorator(
	fn func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	inCh := f.ensure(input)
	outCh := f.ensure(output)

	f.tasks = append(f.tasks, func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	})
}

func (f *Flow) AddMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	allIn := make([]chan string, len(inputs))
	for i, n := range inputs {
		allIn[i] = f.ensure(n)
	}

	out := f.ensure(output)

	f.tasks = append(f.tasks, func(ctx context.Context) error {
		return fn(ctx, allIn, out)
	})
}

func (f *Flow) AddSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	inp := f.ensure(input)

	outList := make([]chan string, len(outputs))
	for i, n := range outputs {
		outList[i] = f.ensure(n)
	}

	f.tasks = append(f.tasks, func(ctx context.Context) error {
		return fn(ctx, inp, outList)
	})
}

func (f *Flow) Run(ctx context.Context) error {
	group, gctx := errgroup.WithContext(ctx)

	for _, t := range f.tasks {
		fn := t

		group.Go(func() error {
			return fn(gctx)
		})
	}

	err := group.Wait()

	f.busGuard.Lock()
	for _, c := range f.bus {
		close(c)
	}
	f.busGuard.Unlock()

	if err != nil {
		return fmt.Errorf("pipeline stopped: %w", err)
	}
	return nil
}

func (f *Flow) Push(chName string, data string) error {
	f.busGuard.RLock()
	ch, ok := f.bus[chName]
	f.busGuard.RUnlock()

	if !ok {
		return ErrChannelMissing
	}

	ch <- data
	return nil
}

func (f *Flow) Pull(name string) (string, error) {
	f.busGuard.RLock()
	ch, ok := f.bus[name]
	f.busGuard.RUnlock()

	if !ok {
		return "", ErrChannelMissing
	}

	v, ok := <-ch
	if !ok {
		return emptyValue, nil
	}

	return v, nil
}
