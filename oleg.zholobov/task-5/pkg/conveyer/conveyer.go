package conveyer

import (
	"context"
	"errors"
	"sync"
)

type decoratorConfig struct {
	fn     func(ctx context.Context, input chan string, output chan string) error
	input  string
	output string
}

type multiplexerConfig struct {
	fn     func(ctx context.Context, inputs []chan string, output chan string) error
	inputs []string
	output string
}

type separatorConfig struct {
	fn      func(ctx context.Context, input chan string, outputs []chan string) error
	input   string
	outputs []string
}

type Conveyer struct {
	size int

	channels map[string]chan string

	decorators   []decoratorConfig
	multiplexers []multiplexerConfig
	separators   []separatorConfig

	mu sync.RWMutex
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
	}
}

func (c *Conveyer) getOrCreateChannel(name string) chan string {
	ch, ok := c.channels[name]
	if !ok {
		ch = make(chan string, c.size)
		c.channels[name] = ch
	}
	return ch
}

func (c *Conveyer) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.getOrCreateChannel(input)
	c.getOrCreateChannel(output)

	c.decorators = append(c.decorators, decoratorConfig{
		fn:     fn,
		input:  input,
		output: output,
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, in := range inputs {
		c.getOrCreateChannel(in)
	}
	c.getOrCreateChannel(output)

	c.multiplexers = append(c.multiplexers, multiplexerConfig{
		fn:     fn,
		inputs: inputs,
		output: output,
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.getOrCreateChannel(input)
	for _, out := range outputs {
		c.getOrCreateChannel(out)
	}

	c.separators = append(c.separators, separatorConfig{
		fn:      fn,
		input:   input,
		outputs: outputs,
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	var wg sync.WaitGroup

	errCh := make(chan error, 1)

	runners := c.buildRunners()

	for _, r := range runners {
		wg.Add(1)
		go func(run runFunc) {
			defer wg.Done()
			if err := run(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(r)
	}

	var finalErr error

	select {
	case <-ctx.Done():
		finalErr = ctx.Err()
	case err := <-errCh:
		finalErr = err
	}

	wg.Wait()

	c.closeAllChannels()

	return finalErr
}

func (c *Conveyer) buildRunners() []runFunc {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var runners []runFunc

	// Decorators
	for _, d := range c.decorators {
		input := c.channels[d.input]
		output := c.channels[d.output]
		fn := d.fn
		runners = append(runners, func(ctx context.Context) error {
			return fn(ctx, input, output)
		})
	}

	// Multiplexers
	for _, m := range c.multiplexers {
		var inputs []chan string
		for _, name := range m.inputs {
			inputs = append(inputs, c.channels[name])
		}
		output := c.channels[m.output]
		fn := m.fn
		runners = append(runners, func(ctx context.Context) error {
			return fn(ctx, inputs, output)
		})
	}

	// Separators
	for _, s := range c.separators {
		input := c.channels[s.input]
		var outputs []chan string
		for _, name := range s.outputs {
			outputs = append(outputs, c.channels[name])
		}
		fn := s.fn
		runners = append(runners, func(ctx context.Context) error {
			return fn(ctx, input, outputs)
		})
	}

	return runners
}

func (c *Conveyer) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for name, ch := range c.channels {
		func() {
			defer func() { _ = recover() }()
			close(ch)
			_ = name
		}()
	}
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.RLock()
	ch, ok := c.channels[input]
	c.mu.RUnlock()

	if !ok {
		return errors.New("chan not found")
	}

	ch <- data
	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, ok := c.channels[output]
	c.mu.RUnlock()

	if !ok {
		return "", errors.New("chan not found")
	}

	val, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return val, nil
}

type runFunc func(ctx context.Context) error
