package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("channel not found")

type conveyer interface {
	RegisterDecorator(fn func(context.Context, chan string, chan string) error, input, output string)
	RegisterMultiplexer(fn func(context.Context, []chan string, chan string) error, inputs []string, output string)
	RegisterSeparator(fn func(context.Context, chan string, []chan string) error, input string, outputs []string)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type pipeline struct {
	bufferSize int
	channels   map[string]chan string
	workers    []func(context.Context) error
	mutex      sync.RWMutex
}

func New(size int) *pipeline {
	return &pipeline{
		bufferSize: size,
		channels:   make(map[string]chan string),
		workers:    make([]func(context.Context) error, 0),
		mutex:      sync.RWMutex{},
	}
}
