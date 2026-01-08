package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChannelNotFound = errors.New("chan not found")

type WorkerFunc func(context.Context) error

type Conveyor struct {
	size     int
	channels map[string]chan string
	workers  []WorkerFunc
	mutex    sync.RWMutex
}

func (conveyor *Conveyor) registerWorker(worker WorkerFunc) {
	conveyor.mutex.Lock()
	conveyor.workers = append(conveyor.workers, worker)
	conveyor.mutex.Unlock()
}

func (conveyor *Conveyor) copyWorkers() []WorkerFunc {
	conveyor.mutex.RLock()
	workers := conveyor.workers
	conveyor.mutex.RUnlock()

	return workers
}

func (conveyor *Conveyor) getOrCreateChannel(name string) chan string {
	conveyor.mutex.Lock()

	channel, ok := conveyor.channels[name]
	if !ok {
		channel = make(chan string, conveyor.size)
		conveyor.channels[name] = channel
	}

	conveyor.mutex.Unlock()

	return channel
}

func (conveyor *Conveyor) getChannel(name string) (chan string, error) {
	conveyor.mutex.RLock()
	channel, ok := conveyor.channels[name]
	conveyor.mutex.RUnlock()

	if !ok {
		return nil, ErrChannelNotFound
	}

	return channel, nil
}

func (conveyor *Conveyor) closeChannels() {
	conveyor.mutex.Lock()

	for _, channel := range conveyor.channels {
		close(channel)
	}

	conveyor.mutex.Unlock()
}

func New(size int) *Conveyor {
	return &Conveyor{
		size:     size,
		channels: make(map[string]chan string),
		workers:  make([]WorkerFunc, 0),
		mutex:    sync.RWMutex{},
	}
}

func (conveyor *Conveyor) RegisterDecorator(
	process func(ctx context.Context, in chan string, out chan string) error,
	inputChannelName string, outputChannelName string,
) {
	inputChannel := conveyor.getOrCreateChannel(inputChannelName)
	outputChannel := conveyor.getOrCreateChannel(outputChannelName)

	worker := func(ctx context.Context) error {
		return process(ctx, inputChannel, outputChannel)
	}
	conveyor.registerWorker(worker)
}

func (conveyor *Conveyor) RegisterMultiplexer(
	process func(ctx context.Context, ins []chan string, out chan string) error,
	inputChannelNames []string, outputChannelName string,
) {
	inputChannels := make([]chan string, len(inputChannelNames))
	for index, name := range inputChannelNames {
		inputChannels[index] = conveyor.getOrCreateChannel(name)
	}

	outputChannel := conveyor.getOrCreateChannel(outputChannelName)

	worker := func(ctx context.Context) error {
		return process(ctx, inputChannels, outputChannel)
	}
	conveyor.registerWorker(worker)
}

func (conveyor *Conveyor) RegisterSeparator(
	processor func(ctx context.Context, in chan string, outs []chan string) error,
	inputChannelName string, outputChannelNames []string,
) {
	inputChannel := conveyor.getOrCreateChannel(inputChannelName)

	outputChannels := make([]chan string, len(outputChannelNames))
	for index, name := range outputChannelNames {
		outputChannels[index] = conveyor.getOrCreateChannel(name)
	}

	worker := func(ctx context.Context) error {
		return processor(ctx, inputChannel, outputChannels)
	}
	conveyor.registerWorker(worker)
}

func (conveyor *Conveyor) Run(ctx context.Context) error {
	group, groupCtx := errgroup.WithContext(ctx)

	workers := conveyor.copyWorkers()

	for _, worker := range workers {
		localWorker := worker

		group.Go(func() error { return localWorker(groupCtx) })
	}

	err := group.Wait()

	conveyor.closeChannels()

	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

func (conveyor *Conveyor) Send(name string, data string) error {
	channel, err := conveyor.getChannel(name)
	if err != nil {
		return err
	}

	defer func() {
		_ = recover()
	}()
	channel <- data

	return nil
}

func (conveyor *Conveyor) Recv(name string) (string, error) {
	channel, err := conveyor.getChannel(name)
	if err != nil {
		return "", err
	}

	data, ok := <-channel
	if !ok {
		return "undefined", nil
	}

	return data, nil
}
