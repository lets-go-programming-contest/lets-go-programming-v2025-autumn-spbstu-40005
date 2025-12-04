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

type Pipeline struct {
	size     int
	channels map[string]chan string
	workers  []func(context.Context) error
	mu       sync.RWMutex
}

func New(size int) *Pipeline {
	return &Pipeline{
		size:     size,
		channels: make(map[string]chan string),
		workers:  make([]func(context.Context) error, 0),
	}
}

func (p *Pipeline) getOrCreateChannel(name string) chan string {
	p.mu.Lock()
	defer p.mu.Unlock()

	channel, exists := p.channels[name]
	if exists {
		return channel
	}

	c := make(chan string, p.size)
	p.channels[name] = c
	return c
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
	inChannels := make([]chan string, len(inputs))

	for i, name := range inputs {
		inChannels[i] = p.getOrCreateChannel(name)
	}

	outCh := p.getOrCreateChannel(output)

	p.workers = append(p.workers, func(ctx context.Context) error {
		return handler(ctx, inChannels, outCh)
	})
}

func (p *Pipeline) RegisterSeparator(
	handler func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	inCh := p.getOrCreateChannel(input)
	outChList := make([]chan string, len(outputs))

	for i, name := range outputs {
		outChList[i] = p.getOrCreateChannel(name)
	}

	p.workers = append(p.workers, func(ctx context.Context) error {
		return handler(ctx, inCh, outChList)
	})
}

func (p *Pipeline) Run(ctx context.Context) error {
	group, derivedCtx := errgroup.WithContext(ctx)

	for _, worker := range p.workers {
		w := worker
		group.Go(func() error {
			return w(derivedCtx)
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

func (p *Pipeline) Send(inputName string, data string) error {
	p.mu.RLock()
	channel, exists := p.channels[inputName]
	p.mu.RUnlock()

	if !exists {
		return ErrChanNotFound
	}

	channel <- data
	return nil
}

func (p *Pipeline) Recv(outputName string) (string, error) {
	p.mu.RLock()
	channel, exists := p.channels[outputName]
	p.mu.RUnlock()

	if !exists {
		return "", ErrChanNotFound
	}

	value, open := <-channel
	if !open {
		return UndefinedValue, nil
	}

	return value, nil
}
