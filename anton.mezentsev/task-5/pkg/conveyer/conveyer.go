package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const undefinedResult = "undefined"

var (
	ErrChannelNotFound  = errors.New("chan not found")
	ErrSourceNamesEmpty = errors.New("sourceNames cannot be empty")
	ErrDestNamesEmpty   = errors.New("destNames cannot be empty")
)

type Task struct {
	execute func(context.Context) error
}

type Pipeline struct {
	mutex         sync.RWMutex
	channels      map[string]chan string
	tasks         []Task
	channelBuffer int
}

func New(size int) *Pipeline {
	return &Pipeline{
		mutex:         sync.RWMutex{},
		channels:      make(map[string]chan string),
		tasks:         make([]Task, 0),
		channelBuffer: size,
	}
}

func (p *Pipeline) getChannel(name string) (chan string, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	channel, exists := p.channels[name]
	if !exists {
		return nil, ErrChannelNotFound
	}

	return channel, nil
}

func (p *Pipeline) Send(inputName string, data string) error {
	channel, err := p.getChannel(inputName)
	if err != nil {
		return err
	}

	channel <- data

	return nil
}

func (p *Pipeline) Recv(outputName string) (string, error) {
	channel, err := p.getChannel(outputName)
	if err != nil {
		return "", err
	}

	data, ok := <-channel
	if !ok {
		return undefinedResult, nil
	}

	return data, nil
}

func (p *Pipeline) getOrCreateChannel(name string) chan string {
	if channel, exists := p.channels[name]; exists {
		return channel
	}

	newChannel := make(chan string, p.channelBuffer)
	p.channels[name] = newChannel

	return newChannel
}

func (p *Pipeline) RegisterDecorator(
	workerFunc func(context.Context, chan string, chan string) error,
	sourceName string,
	destName string,
) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	sourceChannel := p.getOrCreateChannel(sourceName)
	destChannel := p.getOrCreateChannel(destName)

	task := Task{
		execute: func(ctx context.Context) error {
			return workerFunc(ctx, sourceChannel, destChannel)
		},
	}

	p.tasks = append(p.tasks, task)
}

func (p *Pipeline) RegisterMultiplexer(
	workerFunc func(context.Context, []chan string, chan string) error,
	sourceNames []string,
	destName string,
) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(sourceNames) == 0 {
		return ErrSourceNamesEmpty
	}

	sources := make([]chan string, len(sourceNames))
	for i, name := range sourceNames {
		sources[i] = p.getOrCreateChannel(name)
	}

	destChannel := p.getOrCreateChannel(destName)

	task := Task{
		execute: func(ctx context.Context) error {
			return workerFunc(ctx, sources, destChannel)
		},
	}

	p.tasks = append(p.tasks, task)

	return nil
}

func (p *Pipeline) RegisterSeparator(
	workerFunc func(context.Context, chan string, []chan string) error,
	sourceName string,
	destNames []string,
) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(destNames) == 0 {
		return ErrDestNamesEmpty
	}

	sourceChannel := p.getOrCreateChannel(sourceName)

	destinations := make([]chan string, len(destNames))
	for i, name := range destNames {
		destinations[i] = p.getOrCreateChannel(name)
	}

	task := Task{
		execute: func(ctx context.Context) error {
			return workerFunc(ctx, sourceChannel, destinations)
		},
	}

	p.tasks = append(p.tasks, task)

	return nil
}

func (p *Pipeline) closeChannels() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, channel := range p.channels {
		close(channel)
	}
}

func (p *Pipeline) Run(ctx context.Context) error {
	defer p.closeChannels()

	group, groupCtx := errgroup.WithContext(ctx)

	p.mutex.RLock()

	for _, task := range p.tasks {
		currentTask := task

		group.Go(func() error {
			return currentTask.execute(groupCtx)
		})
	}

	p.mutex.RUnlock()

	if err := group.Wait(); err != nil {
		return fmt.Errorf("pipeline execution failed: %w", err)
	}

	return nil
}
