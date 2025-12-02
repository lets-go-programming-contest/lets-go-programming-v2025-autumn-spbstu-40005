package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

type conveyer interface {
	RegisterDecorator(handler func(context.Context, chan string, chan string) error, input, output string)
	RegisterMultiplexer(handler func(context.Context, []chan string, chan string) error, inputs []string, output string)
	RegisterSeparator(handler func(context.Context, chan string, []chan string) error, input string, outputs []string)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

var _ conveyer = (*pipeline)(nil)

type pipeline struct {
	bufferSize int
	channels   map[string]chan string
	workers    []func(context.Context) error
	mutex      sync.RWMutex
	workersMu  sync.RWMutex
}

func New(size int) *pipeline {
	return &pipeline{
		bufferSize: size,
		channels:   make(map[string]chan string),
		workers:    make([]func(context.Context) error, 0),
		mutex:      sync.RWMutex{},
		workersMu:  sync.RWMutex{},
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
	handler func(context.Context, chan string, chan string) error,
	input, output string,
) {
	inChan := p.getOrCreateChannel(input)
	outChan := p.getOrCreateChannel(output)

	worker := func(ctx context.Context) error {
		return handler(ctx, inChan, outChan)
	}

	p.workersMu.Lock()
	defer p.workersMu.Unlock()
	p.workers = append(p.workers, worker)
}

func (p *pipeline) RegisterMultiplexer(
	handler func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	inChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChans[i] = p.getOrCreateChannel(name)
	}

	outChan := p.getOrCreateChannel(output)

	worker := func(ctx context.Context) error {
		return handler(ctx, inChans, outChan)
	}

	p.workersMu.Lock()
	defer p.workersMu.Unlock()
	p.workers = append(p.workers, worker)
}

func (p *pipeline) RegisterSeparator(
	handler func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	inChan := p.getOrCreateChannel(input)

	outChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChans[i] = p.getOrCreateChannel(name)
	}

	worker := func(ctx context.Context) error {
		return handler(ctx, inChan, outChans)
	}

	p.workersMu.Lock()
	defer p.workersMu.Unlock()
	p.workers = append(p.workers, worker)
}

func (p *pipeline) Send(input string, data string) error {
	ch, err := p.getChannel(input)
	if err != nil {
		return err
	}

	ch <- data
	return nil
}

func (p *pipeline) Recv(output string) (string, error) {
	ch, err := p.getChannel(output)
	if err != nil {
		return "", err
	}

	data, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return data, nil
}

func (p *pipeline) closeChannels() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, ch := range p.channels {
		close(ch)
	}
}

func (p *pipeline) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	p.workersMu.RLock()
	workersCopy := make([]func(context.Context) error, len(p.workers))
	copy(workersCopy, p.workers)
	p.workersMu.RUnlock()

	for _, w := range workersCopy {
		workerFunc := w
		group.Go(func() error {
			return workerFunc(ctx)
		})
	}

	err := group.Wait()

	p.closeChannels()

	return err
}
