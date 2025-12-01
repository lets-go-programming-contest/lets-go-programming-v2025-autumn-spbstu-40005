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
	defer c.closeAll()

	errGr, ctx := errgroup.WithContext(ctx)

	for _, worker := range c.workers {
		errGr.Go(func() error {
			return worker(ctx)
		})
	}

	err := errGr.Wait()
	if err != nil {
		return fmt.Errorf("execution failed: %w", err) // Теперь оборачиваем
	}

	return nil
}

func (c *conveyor) closeAll() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, channel := range c.channels {
		close(channel)
	}
}

func (c *conveyor) Send(name string, data string) error {
	c.mutex.RLock()
	channel, exists := c.channels[name]
	c.mutex.RUnlock()

	if !exists {
		return ErrChannelNotFound
	}

	channel <- data
	return nil
}

func (c *conveyor) Recv(name string) (string, error) {
	c.mutex.RLock()
	channel, exists := c.channels[name]
	c.mutex.RUnlock()

	if !exists {
		return "", ErrChannelNotFound
	}

	value, ok := <-channel
	if !ok {
		return undefinedValue, nil
	}

	return value, nil
}

func (c *conveyor) getChannelOrCreate(name string) chan string {
	if channel, ok := c.channels[name]; ok {
		return channel
	}

	channel := make(chan string, c.bufferSize)
	c.channels[name] = channel

	return channel
}
