package conveyer

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

type conveyerImpl struct {
	mu           sync.RWMutex
	channels     map[string]chan string
	decorators   []decoratorConfig
	multiplexers []multiplexerConfig
	separators   []separatorConfig
	size         int
}

type decoratorConfig struct {
	fn     func(ctx context.Context, input chan string, output chan string) error
	input  string
	output string
}

type multiplexerConfig struct {
	fn     func(ctx context.Context, input []chan string, output chan string) error
	inputs []string
	output string
}

type separatorConfig struct {
	fn      func(ctx context.Context, input chan string, output []chan string) error
	input   string
	outputs []string
}

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		channels: make(map[string]chan string),
		size:     size,
	}
}

func (c *conveyerImpl) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *conveyerImpl) getChannel(name string) (chan string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, exists := c.channels[name]
	return ch, exists
}

func (c *conveyerImpl) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.decorators = append(c.decorators, decoratorConfig{
		fn:     fn,
		input:  input,
		output: output,
	})
}

func (c *conveyerImpl) RegisterMultiplexer(
	fn func(ctx context.Context, input []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.multiplexers = append(c.multiplexers, multiplexerConfig{
		fn:     fn,
		inputs: inputs,
		output: output,
	})
}

func (c *conveyerImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, output []chan string) error,
	input string,
	outputs []string,
) {
	c.separators = append(c.separators, separatorConfig{
		fn:      fn,
		input:   input,
		outputs: outputs,
	})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, decorator := range c.decorators {
		decorator := decorator
		inputChan := c.getOrCreateChannel(decorator.input)
		outputChan := c.getOrCreateChannel(decorator.output)

		g.Go(func() error {
			return decorator.fn(ctx, inputChan, outputChan)
		})
	}

	for _, multiplexer := range c.multiplexers {
		multiplexer := multiplexer
		var inputChans []chan string
		for _, input := range multiplexer.inputs {
			inputChans = append(inputChans, c.getOrCreateChannel(input))
		}
		outputChan := c.getOrCreateChannel(multiplexer.output)

		g.Go(func() error {
			return multiplexer.fn(ctx, inputChans, outputChan)
		})
	}

	for _, separator := range c.separators {
		separator := separator
		inputChan := c.getOrCreateChannel(separator.input)
		var outputChans []chan string
		for _, output := range separator.outputs {
			outputChans = append(outputChans, c.getOrCreateChannel(output))
		}

		g.Go(func() error {
			return separator.fn(ctx, inputChan, outputChans)
		})
	}

	err := g.Wait()

	c.mu.Lock()
	for name, ch := range c.channels {
		close(ch)
		delete(c.channels, name)
	}
	c.mu.Unlock()

	return err
}

func (c *conveyerImpl) Send(input string, data string) error {
	ch := c.getOrCreateChannel(input)

	select {
	case ch <- data:
		return nil
	default:
		return fmt.Errorf("channel is full")
	}
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	ch := c.getOrCreateChannel(output)

	select {
	case data, ok := <-ch:
		if !ok {
			return "undefined", nil
		}
		return data, nil
	default:
		return "", fmt.Errorf("no data available")
	}
}
