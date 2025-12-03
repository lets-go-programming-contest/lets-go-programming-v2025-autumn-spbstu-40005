package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const undefinedValue = "undefined"

var ErrChannelNotFound = errors.New("chan not found")

type conveyor struct {
	bufferSize int
	channels   map[string]chan string
	workers    []func(ctx context.Context) error
	mutex      sync.RWMutex
}

func New(size int) *conveyor {
	return &conveyor{
		bufferSize: size,
		channels:   make(map[string]chan string),
		workers:    []func(ctx context.Context) error{},
		mutex:      sync.RWMutex{},
	}
}

func (c *conveyor) getChannel(name string) (chan string, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	channel, exists := c.channels[name]
	return channel, exists
}

func (c *conveyor) getChannelOrCreate(name string) chan string {
	c.mutex.RLock()
	channel, exists := c.channels[name]
	c.mutex.RUnlock()

	if exists {
		return channel
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if channel, exists := c.channels[name]; exists {
		return channel
	}

	channel = make(chan string, c.bufferSize)
	c.channels[name] = channel
	return channel
}

func (c *conveyor) RegisterDecorator(
	processor func(ctx context.Context, input chan string, output chan string) error,
	inputName,
	outputName string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inputChannel := c.getChannelOrCreate(inputName)
	outputChannel := c.getChannelOrCreate(outputName)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return processor(ctx, inputChannel, outputChannel)
	})
}

func (c *conveyor) RegisterMultiplexer(
	processor func(ctx context.Context, inputs []chan string, output chan string) error,
	inputNames []string,
	outputName string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	outputChannel := c.getChannelOrCreate(outputName)
	inputChannels := make([]chan string, len(inputNames))

	for i, name := range inputNames {
		inputChannels[i] = c.getChannelOrCreate(name)
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return processor(ctx, inputChannels, outputChannel)
	})
}

func (c *conveyor) RegisterSeparator(
	processor func(ctx context.Context, input chan string, outputs []chan string) error,
	inputName string,
	outputNames []string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inputChannel := c.getChannelOrCreate(inputName)
	outputChannels := make([]chan string, len(outputNames))

	for i, name := range outputNames {
		outputChannels[i] = c.getChannelOrCreate(name)
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return processor(ctx, inputChannel, outputChannels)
	})
}

func (c *conveyor) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errGr, ctx := errgroup.WithContext(ctx)

	c.mutex.RLock()
	workersCopy := make([]func(context.Context) error, len(c.workers))
	copy(workersCopy, c.workers)
	c.mutex.RUnlock()

	for _, worker := range workersCopy {
		worker := worker
		errGr.Go(func() error {
			return worker(ctx)
		})
	}

	err := errGr.Wait()

	c.closeAll()

	if err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

func (c *conveyor) closeAll() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for name, channel := range c.channels {
		select {
		case <-channel:
		default:
		}
		close(channel)
		delete(c.channels, name)
	}
}

func (c *conveyor) Send(name string, data string) error {
	channel, exists := c.getChannel(name)
	if !exists {
		return ErrChannelNotFound
	}

	select {
	case channel <- data:
		return nil
	default:
		return errors.New("channel is full")
	}
}

func (c *conveyor) Recv(name string) (string, error) {
	channel, exists := c.getChannel(name)
	if !exists {
		return "", ErrChannelNotFound
	}

	select {
	case value, ok := <-channel:
		if !ok {
			return undefinedValue, nil
		}
		return value, nil
	default:
		return "", errors.New("no data available")
	}
}
