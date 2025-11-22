package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound   = errors.New("chan not found")
	ErrChannelFull    = errors.New("channel is full")
	ErrNoDataAvailable = errors.New("no data available")
	ErrUndefined      = "undefined"
)

type DecoratorFunc func(ctx context.Context, input chan string, output chan string) error
type MultiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error
type SeparatorFunc func(ctx context.Context, input chan string, outputs []chan string) error

type handler struct {
	fn      interface{}
	inputs  []string
	outputs []string
}

type Conveyer struct {
	size      int
	channels  map[string]chan string
	mu        sync.RWMutex
	handlers  []handler
	cancel    context.CancelFunc
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]handler, 0),
	}
}

func (c *Conveyer) getOrCreateChannel(id string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[id]; exists {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[id] = ch
	return ch
}

func (c *Conveyer) getChannel(id string) (chan string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if ch, exists := c.channels[id]; exists {
		return ch, nil
	}

	return nil, ErrChanNotFound
}

func (c *Conveyer) RegisterDecorator(
	fn DecoratorFunc,
	input string,
	output string,
) {
	c.handlers = append(c.handlers, handler{
		fn:      fn,
		inputs:  []string{input},
		outputs: []string{output},
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn MultiplexerFunc,
	inputs []string,
	output string,
) {
	c.handlers = append(c.handlers, handler{
		fn:      fn,
		inputs:  inputs,
		outputs: []string{output},
	})
}

func (c *Conveyer) RegisterSeparator(
	fn SeparatorFunc,
	input string,
	outputs []string,
) {
	c.handlers = append(c.handlers, handler{
		fn:      fn,
		inputs:  []string{input},
		outputs: outputs,
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, c.cancel = context.WithCancel(ctx)
	defer c.cancel()

	g, ctx := errgroup.WithContext(ctx)

	for _, h := range c.handlers {
		h := h
		g.Go(func() error {
			switch fn := h.fn.(type) {
			case DecoratorFunc:
				inputChan := c.getOrCreateChannel(h.inputs[0])
				outputChan := c.getOrCreateChannel(h.outputs[0])
				return fn(ctx, inputChan, outputChan)
			case MultiplexerFunc:
				inputChans := make([]chan string, len(h.inputs))
				for i, input := range h.inputs {
					inputChans[i] = c.getOrCreateChannel(input)
				}
				outputChan := c.getOrCreateChannel(h.outputs[0])
				return fn(ctx, inputChans, outputChan)
			case SeparatorFunc:
				inputChan := c.getOrCreateChannel(h.inputs[0])
				outputChans := make([]chan string, len(h.outputs))
				for i, output := range h.outputs {
					outputChans[i] = c.getOrCreateChannel(output)
				}
				return fn(ctx, inputChan, outputChans)
			default:
				return errors.New("unknown handler type")
			}
		})
	}

	return g.Wait()
}

func (c *Conveyer) Send(input string, data string) error {
	ch, err := c.getChannel(input)
	if err != nil {
		return ErrChanNotFound
	}

	select {
	case ch <- data:
		return nil
	default:
		return ErrChannelFull
	}
}

func (c *Conveyer) Recv(output string) (string, error) {
	ch, err := c.getChannel(output)
	if err != nil {
		return "", ErrChanNotFound
	}

	select {
	case data, ok := <-ch:
		if !ok {
			return ErrUndefined, nil
		}
		return data, nil
	default:
		return "", ErrNoDataAvailable
	}
}

func (c *Conveyer) Stop() {
	if c.cancel != nil {
		c.cancel()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for id, ch := range c.channels {
		close(ch)
		delete(c.channels, id)
	}
}
