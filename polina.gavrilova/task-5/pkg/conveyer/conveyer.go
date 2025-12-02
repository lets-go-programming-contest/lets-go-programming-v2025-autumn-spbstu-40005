package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("channel not found")

type conveyer interface {
	RegisterDecorator(fn func(context.Context, chan string, chan string) error, input, output string)
	RegisterMultiplexer(fn func(context.Context, []chan string, chan string) error, inputs []string, output string)
	RegisterSeparator(fn func(context.Context, chan string, []chan string) error, input string, outputs []string)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type pipeline struct {
	bufferSize int
	channels   map[string]chan string
	workers    []func(context.Context) error
	mutex      sync.RWMutex
}

func New(size int) *pipeline {
	return &pipeline{
		bufferSize: size,
		channels:   make(map[string]chan string),
		workers:    make([]func(context.Context) error, 0),
		mutex:      sync.RWMutex{},
	}
}

func (p *pipeline) getOrCreateChannel(name string) chan string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if ch, ok := p.channels[name]; ok {
		return ch
	}

	newChan := make(chan string, p.bufferSize)
	p.channels[name] = newChan

	return newChan
}

func (p *pipeline) getChannel(name string) (chan string, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if ch, ok := p.channels[name]; ok {
		return ch, nil
	}

	return nil, ErrChanNotFound
}

func (p *pipeline) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	input, output string,
) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	inChan := p.getOrCreateChannel(input)
	outChan := p.getOrCreateChannel(output)

	worker := func(ctx context.Context) error {
		return fn(ctx, inChan, outChan)
	}

	p.workers = append(p.workers, worker)
}

func (p *pipeline) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	inChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChans[i] = p.getOrCreateChannel(name)
	}

	outChan := p.getOrCreateChannel(output)

	worker := func(ctx context.Context) error {
		return fn(ctx, inChans, outChan)
	}

	p.workers = append(p.workers, worker)
}

func (p *pipeline) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	inChan := p.getOrCreateChannel(input)
	outChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChans[i] = p.getOrCreateChannel(name)
	}

	worker := func(ctx context.Context) error {
		return fn(ctx, inChan, outChans)
	}

	p.workers = append(p.workers, worker)
}
