package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

const UndefinedValue = "undefined"

type ConveyerInterface interface {
	RegisterDecorator(
		handler func(context.Context, chan string, chan string) error,
		input string,
		output string,
	)

	RegisterMultiplexer(
		handler func(context.Context, []chan string, chan string) error,
		inputs []string,
		output string,
	)

	RegisterSeparator(
		handler func(context.Context, chan string, []chan string) error,
		input string,
		outputs []string,
	)

	Run(ctx context.Context) error
	Send(inputName string, data string) error
	Recv(outputName string) (string, error)
}

type WorkerFunc func(context.Context) error

type Pipeline struct {
	size     int
	channels map[string]chan string
	workers  []WorkerFunc
	mu       sync.RWMutex
}

func New(size int) *Pipeline {
	return &Pipeline{
		size:     size,
		channels: make(map[string]chan string),
		workers:  make([]WorkerFunc, 0),
		mu:       sync.RWMutex{},
	}
}

func (p *Pipeline) getOrCreateChannel(name string) chan string {
	p.mu.Lock()
	defer p.mu.Unlock()

	existing, exists := p.channels[name]
	if exists {
		return existing
	}

	ch := make(chan string, p.size)
	p.channels[name] = ch

	return ch
}

func (p *Pipeline) RegisterDecorator(
	handler func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	inCh := p.getOrCreateChannel(input)
	outCh := p.getOrCreateChannel(output)

	p.workers = append(p.workers, func(ctx context.Context) error {
		return handler(ctx, inCh, outCh)
	})
}

func (p *Pipeline) RegisterMultiplexer(
	handler func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	inChList := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChList[i] = p.getOrCreateChannel(name)
	}

	outCh := p.getOrCreateChannel(output)

	p.workers = append(p.workers, func(ctx context.Context) error {
		return handler(ctx, inChList, outCh)
	})
}

func (p *Pipeline) RegisterSeparator(
	handler func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	inCh := p.getOrCreateChannel(input)

	outList := make([]chan string, len(outputs))
	for i, name := range outputs {
		outList[i] = p.getOrCreateChannel(name)
	}

	p.workers = append(p.workers, func(ctx context.Context) error {
		return handler(ctx, inCh, outList)
	})
}

func (p *Pipeline) Run(ctx context.Context) error {
	group, derivedCtx := errgroup.WithContext(ctx)

	for _, job := range p.workers {
		task := job

		group.Go(func() error {
			return task(derivedCtx)
		})
	}

	err := group.Wait()

	p.mu.Lock()
	for _, ch := range p.channels {
		close(ch)
	}
	p.mu.Unlock()

	if err != nil {
		return fmt.Errorf("conveyer run failed: %w", err)
	}

	return nil
}

func (p *Pipeline) lookupChannel(name string) (chan string, bool) {
	p.mu.RLock()
	ch, ok := p.channels[name]
	p.mu.RUnlock()

	return ch, ok
}

func (p *Pipeline) Send(inputName string, data string) error {
	channel, exists := p.lookupChannel(inputName)
	if !exists {
		return ErrChanNotFound
	}
	channel <- data

	return nil
}

func (p *Pipeline) Recv(outputName string) (string, error) {
	channel, exists := p.lookupChannel(outputName)
	if !exists {
		return "", ErrChanNotFound
	}
	value, open := <-channel

	if !open {

		return UndefinedValue, nil
	}

	return value, nil
}
