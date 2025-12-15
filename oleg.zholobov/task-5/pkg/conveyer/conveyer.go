package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const undefinedValue = "undefined"

var errMsgChannelNotFound = errors.New("chan not found")

type WorkerFunc func(ctx context.Context) error

type Conveyer struct {
	size     int
	channels map[string]chan string
	workers  []WorkerFunc
	mutex    sync.RWMutex
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		workers:  make([]WorkerFunc, 0),
		mutex:    sync.RWMutex{},
	}
}

func (c *Conveyer) getOrCreateChannel(name string) chan string {
	if ch, ok := c.channels[name]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch

	return ch
}

func (c *Conveyer) getChannel(name string) (chan string, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	channel, exists := c.channels[name]
	if !exists {
		return nil, errMsgChannelNotFound
	}
	return channel, nil
}

func (c *Conveyer) RegisterDecorator(
	callback func(ctx context.Context, in chan string, out chan string) error,
	input string,
	output string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inputChannel := c.getOrCreateChannel(input)
	outputChannel := c.getOrCreateChannel(output)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return callback(ctx, inputChannel, outputChannel)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	callback func(ctx context.Context, ins []chan string, out chan string) error,
	inputs []string,
	output string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	outputChannel := c.getOrCreateChannel(output)
	inputChannels := make([]chan string, len(inputs))

	for i, name := range inputs {
		inputChannels[i] = c.getOrCreateChannel(name)
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return callback(ctx, inputChannels, outputChannel)
	})
}

func (c *Conveyer) RegisterSeparator(
	callback func(ctx context.Context, in chan string, outs []chan string) error,
	input string,
	output []string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	inputChannel := c.getOrCreateChannel(input)
	outputChannels := make([]chan string, len(output))

	for i, name := range output {
		outputChannels[i] = c.getOrCreateChannel(name)
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return callback(ctx, inputChannel, outputChannels)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	defer c.closeAll()

	errGroup, ctx := errgroup.WithContext(ctx)

	for _, w := range c.workers {
		errGroup.Go(func() error {
			return w(ctx)
		})
	}

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("execution failed: %w", err)
	}

	return nil
}

func (c *Conveyer) closeAll() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, ch := range c.channels {
		close(ch)
	}
}

func (c *Conveyer) Send(name string, data string) error {
	channel, err := c.getChannel(name)
	if err != nil {
		return err
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(name string) (string, error) {
	channel, err := c.getChannel(name)
	if err != nil {
		return "", err
	}

	value, ok := <-channel
	if !ok {
		return undefinedValue, nil
	}

	return value, nil
}
