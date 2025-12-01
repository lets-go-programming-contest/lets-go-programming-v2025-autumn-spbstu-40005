package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const undefinedResult = "undefined"

var ErrChannelNotFound = errors.New("chan not found")

type Pipeline struct {
	mutex         sync.RWMutex
	channels      map[string]chan string
	tasks         []func(context.Context) error
	channelBuffer int
}

func New(size int) *Pipeline {
	return &Pipeline{
		mutex:         sync.RWMutex{},
		channels:      make(map[string]chan string),
		tasks:         make([]func(context.Context) error, 0),
		channelBuffer: size,
	}
}

func (p *Pipeline) Send(inputName string, data string) error {
	p.mutex.RLock()
	channel, exists := p.channels[inputName]
	p.mutex.RUnlock()

	if !exists {
		return ErrChannelNotFound
	}

	channel <- data

	return nil
}

func (p *Pipeline) Recv(outputName string) (string, error) {
	p.mutex.RLock()
	channel, exists := p.channels[outputName]
	p.mutex.RUnlock()

	if !exists {
		return "", ErrChannelNotFound
	}

	data, ok := <-channel
	if !ok {
		return undefinedResult, nil
	}

	return data, nil
}

func (p *Pipeline) ensureChannelExists(name string) chan string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if ch, exists := p.channels[name]; exists {
		return ch
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
	sourceChan := p.ensureChannelExists(sourceName)
	destChan := p.ensureChannelExists(destName)

	task := func(ctx context.Context) error {
		return workerFunc(ctx, sourceChan, destChan)
	}

	p.tasks = append(p.tasks, task)
}

func (p *Pipeline) RegisterMultiplexer(
	workerFunc func(context.Context, []chan string, chan string) error,
	sourceNames []string,
	destName string,
) {
	sources := make([]chan string, len(sourceNames))
	for i, name := range sourceNames {
		sources[i] = p.ensureChannelExists(name)
	}

	destChan := p.ensureChannelExists(destName)

	task := func(ctx context.Context) error {
		return workerFunc(ctx, sources, destChan)
	}

	p.tasks = append(p.tasks, task)
}

func (p *Pipeline) RegisterSeparator(
	workerFunc func(context.Context, chan string, []chan string) error,
	sourceName string,
	destNames []string,
) {
	sourceChan := p.ensureChannelExists(sourceName)

	destinations := make([]chan string, len(destNames))
	for i, name := range destNames {
		destinations[i] = p.ensureChannelExists(name)
	}

	task := func(ctx context.Context) error {
		return workerFunc(ctx, sourceChan, destinations)
	}

	p.tasks = append(p.tasks, task)
}

func (p *Pipeline) Run(ctx context.Context) error {
	group, groupCtx := errgroup.WithContext(ctx)

	for _, task := range p.tasks {
		currentTask := task

		group.Go(func() error {
			return currentTask(groupCtx)
		})
	}

	err := group.Wait()

	p.mutex.Lock()
	for _, ch := range p.channels {
		close(ch)
	}

	p.mutex.Unlock()

	if err != nil {
		return fmt.Errorf("pipeline execution failed: %w", err)
	}

	return nil
}
