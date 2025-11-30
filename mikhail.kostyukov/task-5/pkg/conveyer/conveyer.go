package conveyer

import (
	"context"
	"errors"
	"sync"
)

const UndefinedResult = "undefined"

var ErrChannelNotFound = errors.New("chan not found")

type Pipeline struct {
	mutex sync.RWMutex
	channels map[string]chan string
	tasks []func(context.Context) error
	channelBuffer int
}

func New(size int) *Pipeline {
	return &Pipeline{
		channels: make(map[string]chan string),
		tasks: make([]func(context.Context) error, 0),
		channelBuffer: size,
	}
}

func (p *Pipeline) Send(inputName string, data string) error {
	p.mutex.RLock()
	ch, exists := p.channels[inputName]
	p.mutex.RUnlock()

	if !exists {
		return ErrChannelNotFound
	}

	ch <- data

	return nil
}

func (p *Pipeline) Recv(outputName string) (string, error) {
	p.mutex.RLock()
	ch, exists := p.channels[outputName]
	p.mutex.RUnlock()

	if !exists {
		return "", ErrChannelNotFound
	}

	data, ok := <- ch
	if !ok {
		return UndefinedResult, nil
	}

	return data, nil
}
