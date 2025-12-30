package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

const undefined = "undefined"

type TaskFunc func(ctx context.Context) error

type Conveyer struct {
	mu       sync.RWMutex
	channels map[string]chan string
	size     int
	tasks    []TaskFunc
}

func New(size int) *Conveyer {
	return &Conveyer{
		mu:       sync.RWMutex{},
		channels: make(map[string]chan string),
		size:     size,
		tasks:    []TaskFunc{},
	}
}

func (c *Conveyer) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if channel, exists := c.channels[name]; exists {
		return channel
	}

	channel := make(chan string, c.size)
	c.channels[name] = channel

	return channel
}

func (c *Conveyer) RegisterDecorator(
	decoratorFunc func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inputChan := c.getOrCreateChannel(input)
	outputChan := c.getOrCreateChannel(output)

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return decoratorFunc(ctx, inputChan, outputChan)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	multiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputChannels := make([]chan string, len(inputs))
	for i, inputName := range inputs {
		inputChannels[i] = c.getOrCreateChannel(inputName)
	}

	outputChannel := c.getOrCreateChannel(output)

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return multiplexerFunc(ctx, inputChannels, outputChannel)
	})
}

func (c *Conveyer) RegisterSeparator(
	separatorFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inputChannel := c.getOrCreateChannel(input)

	outputChannels := make([]chan string, len(outputs))
	for i, outputName := range outputs {
		outputChannels[i] = c.getOrCreateChannel(outputName)
	}

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return separatorFunc(ctx, inputChannel, outputChannels)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	waitGroup := sync.WaitGroup{}
	errorChannel := make(chan error, len(c.tasks))

	waitGroup.Add(len(c.tasks))

	for _, currentTask := range c.tasks {
		taskCopy := currentTask
		taskFunc := func() {
			defer waitGroup.Done()

			if err := taskCopy(ctx); err != nil {
				select {
				case errorChannel <- err:
				default:
				}
			}
		}

		go taskFunc()
	}

	go func() {
		waitGroup.Wait()
		close(errorChannel)

		c.mu.Lock()
		for _, channel := range c.channels {
			close(channel)
		}

		c.mu.Unlock()
	}()

	select {
	case err := <-errorChannel:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.RLock()
	channel, exists := c.channels[input]
	c.mu.RUnlock()

	if !exists {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	channel, exists := c.channels[output]
	c.mu.RUnlock()

	if !exists {
		return "", ErrChanNotFound
	}

	data, ok := <-channel
	if !ok {
		return undefined, nil
	}

	return data, nil
}
