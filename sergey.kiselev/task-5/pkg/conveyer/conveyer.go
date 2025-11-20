package conveyer

import (
	"context"
	"errors"
)

type conveyer interface {
	RegisterDecorator(
		fn func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		fn func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		fn func(ctx context.Context, input chan string, outputs []chan string) error,
		input string,
		outputs []string,
	)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type conveyerImpl struct {
	size     int
	channels map[string]chan string
}

func New(size int) conveyer {
	return &conveyerImpl{
		size:     size,
		channels: make(map[string]chan string),
	}
}

func (c *conveyerImpl) RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error, input, output string) {
}

func (c *conveyerImpl) RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string) {
}

func (c *conveyerImpl) RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string) {
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	return nil
}

func (c *conveyerImpl) Send(input string, data string) error {
	return errors.New("chan not found")
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	return "", errors.New("chan not found")
}
