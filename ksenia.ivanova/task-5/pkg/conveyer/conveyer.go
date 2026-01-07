package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChannelFull     = errors.New("chan is full")
	ErrChannelNotFound = errors.New("chan not found")
)

type WorkerFunc func(context.Context) error

type Conveyor struct {
	size     int
	channels map[string]chan string
	workers  []WorkerFunc
	mutex    sync.RWMutex
}

func (conveyer *Conveyor) registerWorker(worker WorkerFunc) {
	conveyer.mutex.Lock()
	defer conveyer.mutex.Unlock()

	conveyer.workers = append(conveyer.workers, worker)
}

func (conveyer *Conveyor) getOrCreateChannel(name string) chan string {
	conveyer.mutex.RLock()
	if channel, ok := conveyer.channels[name]; ok {
		conveyer.mutex.RUnlock()
		return channel
	}
	conveyer.mutex.RUnlock()

	conveyer.mutex.Lock()
	defer conveyer.mutex.Unlock()

	if channel, ok := conveyer.channels[name]; ok {
		return channel
	}

	channel := make(chan string, conveyer.size)
	conveyer.channels[name] = channel

	return channel
}

func (conveyer *Conveyor) getChannel(name string) (chan string, error) {
	conveyer.mutex.RLock()
	defer conveyer.mutex.RUnlock()

	channel, ok := conveyer.channels[name]
	if !ok {
		return nil, ErrChannelNotFound
	}

	return channel, nil
}

func (conveyer *Conveyor) closeChannels() {
	conveyer.mutex.Lock()
	defer conveyer.mutex.Unlock()

	for _, channel := range conveyer.channels {
		close(channel)
	}
}

func New(size int) *Conveyor {
	return &Conveyor{
		size:     size,
		channels: make(map[string]chan string),
		workers:  make([]WorkerFunc, 0),
		mutex:    sync.RWMutex{},
	}
}

func (conveyer *Conveyor) RegisterDecorator(
	processor func(ctx context.Context, in chan string, out chan string) error,
	inName string, outName string,
) {
	inChan := conveyer.getOrCreateChannel(inName)
	outChan := conveyer.getOrCreateChannel(outName)

	conveyer.registerWorker(func(ctx context.Context) error {
		return processor(ctx, inChan, outChan)
	})
}

func (conveyer *Conveyor) RegisterMultiplexer(
	processor func(ctx context.Context, ins []chan string, out chan string) error,
	inNames []string, outName string,
) {
	inChans := make([]chan string, len(inNames))
	for idx, name := range inNames {
		inChans[idx] = conveyer.getOrCreateChannel(name)
	}

	outChan := conveyer.getOrCreateChannel(outName)

	conveyer.registerWorker(func(ctx context.Context) error {
		return processor(ctx, inChans, outChan)
	})
}

func (conveyer *Conveyor) RegisterSeparator(
	processor func(ctx context.Context, in chan string, outs []chan string) error,
	inName string, outNames []string,
) {
	inChan := conveyer.getOrCreateChannel(inName)

	outChans := make([]chan string, len(outNames))
	for i, name := range outNames {
		outChans[i] = conveyer.getOrCreateChannel(name)
	}

	conveyer.registerWorker(func(ctx context.Context) error {
		return processor(ctx, inChan, outChans)
	})
}

func (conveyer *Conveyor) Run(parentCtx context.Context) error {
	defer conveyer.closeChannels()

	errGroup, ctx := errgroup.WithContext(parentCtx)

	conveyer.mutex.RLock()

	for _, worker := range conveyer.workers {
		localWorker := worker

		errGroup.Go(func() error { return localWorker(ctx) })
	}

	conveyer.mutex.RUnlock()

	err := errGroup.Wait()
	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

func (conveyer *Conveyor) Send(name string, data string) error {
	channel, err := conveyer.getChannel(name)
	if err != nil {
		return err
	}

	select {
	case channel <- data:
		return nil
	default:
		return ErrChannelFull
	}
}

func (conveyer *Conveyor) Recv(name string) (string, error) {
	channel, err := conveyer.getChannel(name)
	if err != nil {
		return "", err
	}

	data, ok := <-channel
	if !ok {
		return "no data", nil
	}

	return data, nil
}
