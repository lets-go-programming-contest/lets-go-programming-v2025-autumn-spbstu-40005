package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChannel = errors.New("chan not found")

const undefinedValue = "undefined"

type Handlers func(ctx context.Context) error

type ConveyerConfig struct {
	mutex    sync.RWMutex
	channels map[string]chan string
	size     int
	handlers []Handlers
}

type conveyerImpl struct {
	config ConveyerConfig
}

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		config: ConveyerConfig{
			mutex:    sync.RWMutex{},
			channels: make(map[string]chan string),
			size:     size,
			handlers: make([]Handlers, 0),
		},
	}
}

func (c *conveyerImpl) getChannelOrCreate(name string) chan string {
	cfg := &c.config
	if ch, exists := cfg.channels[name]; exists {
		return ch
	}

	ch := make(chan string, cfg.size)
	cfg.channels[name] = ch

	return ch
}

func (c *conveyerImpl) getChannel(name string) (chan string, error) {
	cfg := &c.config
	cfg.mutex.RLock()
	defer cfg.mutex.RUnlock()

	ch, exists := cfg.channels[name]
	if !exists {
		return nil, fmt.Errorf("%w: channel '%s' not found", ErrChannel, name)
	}

	return ch, nil
}

func (c *conveyerImpl) RegisterDecorator(
	function func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	cfg := &c.config
	cfg.mutex.Lock()
	defer cfg.mutex.Unlock()

	inputChan := c.getChannelOrCreate(input)
	outputChan := c.getChannelOrCreate(output)

	cfg.handlers = append(cfg.handlers, func(ctx context.Context) error {
		return function(ctx, inputChan, outputChan)
	})
}

func (c *conveyerImpl) RegisterMultiplexer(
	function func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	cfg := &c.config
	cfg.mutex.Lock()
	defer cfg.mutex.Unlock()

	outputChan := c.getChannelOrCreate(output)
	inputChans := make([]chan string, len(inputs))

	for i, input := range inputs {
		inputChans[i] = c.getChannelOrCreate(input)
	}

	cfg.handlers = append(cfg.handlers, func(ctx context.Context) error {
		return function(ctx, inputChans, outputChan)
	})
}

func (c *conveyerImpl) RegisterSeparator(
	function func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	cfg := &c.config
	cfg.mutex.Lock()
	defer cfg.mutex.Unlock()

	inputChan := c.getChannelOrCreate(input)
	outputChans := make([]chan string, len(outputs))

	for i, output := range outputs {
		outputChans[i] = c.getChannelOrCreate(output)
	}

	cfg.handlers = append(cfg.handlers, func(ctx context.Context) error {
		return function(ctx, inputChan, outputChans)
	})
}

func (c *conveyerImpl) closeAllChannels() {
	cfg := &c.config
	cfg.mutex.Lock()
	defer cfg.mutex.Unlock()

	for _, ch := range cfg.channels {
		close(ch)
	}
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	cfg := &c.config
	group, ctx := errgroup.WithContext(ctx)

	for _, h := range cfg.handlers {
		group.Go(func() error {
			return h(ctx)
		})
	}

	err := group.Wait()

	c.closeAllChannels()

	if err != nil {
		return fmt.Errorf("conveyer execution failed: %w", err)
	}

	return nil
}

func (c *conveyerImpl) Send(input string, data string) error {
	ch, err := c.getChannel(input)
	if err != nil {
		return err
	}

	ch <- data

	return nil
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	ch, err := c.getChannel(output)
	if err != nil {
		return "", err
	}

	data, ok := <-ch
	if !ok {
		return undefinedValue, nil
	}

	return data, nil
}
