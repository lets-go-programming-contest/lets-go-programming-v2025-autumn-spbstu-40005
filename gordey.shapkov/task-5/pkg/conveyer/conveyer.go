package conveyer

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

type Pipeline struct {
	size     int
	channels map[string]chan string
	handlers []func(ctx context.Context) error
}

func New(size int) *Pipeline {
	return &Pipeline{
		size:     size,
		channels: make(map[string]chan string),
		handlers: []func(ctx context.Context) error{},
	}
}

func (pipe *Pipeline) register(ch string) chan string {
	if _, exists := pipe.channels[ch]; !exists {
		pipe.channels[ch] = make(chan string, pipe.size)
	}

	return pipe.channels[ch]
}

func (pipe *Pipeline) RegisterDecorator(
	function func(
		ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string,
	output string,
) {
	in := pipe.register(input)
	out := pipe.register(output)
	pipe.handlers = append(pipe.handlers, func(ctx context.Context) error {
		return function(ctx, in, out)
	})
}

func (pipe *Pipeline) RegisterMultiplexer(
	function func(
		ctx context.Context,
		inputs []chan string,
		output chan string,
	) error,
	inputs []string,
	output string,
) {
	ins := make([]chan string, len(inputs))
	for i, ch := range inputs {
		ins[i] = pipe.register(ch)
	}

	out := pipe.register(output)
	pipe.handlers = append(pipe.handlers, func(ctx context.Context) error {
		return function(ctx, ins, out)
	})
}
func (pipe *Pipeline) RegisterSeparator(
	function func(
		ctx context.Context,
		input chan string,
		outputs []chan string,
	) error,
	input string,
	outputs []string,
) {
	inChan := pipe.register(input)
	outs := make([]chan string, len(outputs))

	for i, ch := range outputs {
		outs[i] = pipe.register(ch)
	}

	pipe.handlers = append(pipe.handlers, func(ctx context.Context) error {
		return function(ctx, inChan, outs)
	})
}

func (pipe *Pipeline) Run(ctx context.Context) error {
	errgr, ctx := errgroup.WithContext(ctx)

	for _, handler := range pipe.handlers {
		errgr.Go(func() error {
			return handler(ctx)
		})
	}

	err := errgr.Wait()

	for _, ch := range pipe.channels {
		close(ch)
	}

	return err
}

func (pipe *Pipeline) Send(input string, data string) error {
	ch, exists := pipe.channels[input]
	if !exists {
		return ErrChanNotFound
	}

	ch <- data

	return nil
}

func (pipe *Pipeline) Recv(output string) (string, error) {
	ch, exists := pipe.channels[output]
	if !exists {
		return "", ErrChanNotFound
	}

	data, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return data, nil
}
