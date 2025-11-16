package conveyer

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	size     int
	channels map[string]chan string
	handlers []func(ctx context.Context) error
}

func New(size int) *Conveyer {
	if size < 0 {
		panic("invalid chan size")
	}

	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		handlers: []func(ctx context.Context) error{},
	}
}

func (p *Conveyer) register(name string) chan string {
	if ch, ok := p.channels[name]; ok {
		return ch
	}

	ch := make(chan string, p.size)
	p.channels[name] = ch

	return ch
}

func (p *Conveyer) RegisterDecorator(
	fn func( //nolint:varnamelen
		ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string,
	output string,
) {
	p.handlers = append(p.handlers, func(ctx context.Context) error {
		return fn(ctx, p.register(input), p.register(output))
	})
}

func (p *Conveyer) RegisterMultiplexer(
	fn func( //nolint:varnamelen
		ctx context.Context,
		inputs []chan string,
		output chan string,
	) error,
	inputs []string,
	output string,
) {
	if len(inputs) == 0 {
		panic("empty inputs")
	}

	inChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChans[i] = p.register(name)
	}

	p.handlers = append(p.handlers, func(ctx context.Context) error {
		return fn(ctx, inChans, p.register(output))
	})
}

func (p *Conveyer) RegisterSeparator(
	fn func( //nolint:varnamelen
		ctx context.Context,
		input chan string,
		outputs []chan string,
	) error,
	input string,
	outputs []string,
) {
	if len(outputs) == 0 {
		panic("empty outputs")
	}

	outChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChans[i] = p.register(name)
	}

	p.handlers = append(p.handlers, func(ctx context.Context) error {
		return fn(ctx, p.register(input), outChans)
	})
}

func (p *Conveyer) Run(ctx context.Context) error {
	defer func() {
		for _, ch := range p.channels {
			close(ch)
		}
	}()

	errgr, ctx := errgroup.WithContext(ctx)

	for _, handler := range p.handlers {
		errgr.Go(func() error {
			return handler(ctx)
		})
	}

	if err := errgr.Wait(); err != nil {
		return fmt.Errorf("run conveyer: %w", err)
	}

	return nil
}

func (p *Conveyer) Send(input string, data string) error {
	ch, ok := p.channels[input]
	if !ok {
		return ErrChanNotFound
	}

	ch <- data

	return nil
}

func (p *Conveyer) Recv(output string) (string, error) {
	ch, ok := p.channels[output] //nolint:varnamelen
	if !ok {
		return "", ErrChanNotFound
	}

	data, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return data, nil
}
