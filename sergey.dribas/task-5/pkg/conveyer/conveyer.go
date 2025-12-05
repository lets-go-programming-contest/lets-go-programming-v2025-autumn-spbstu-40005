package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrChanNotFound     = errors.New("chan not found")
	ErrSendFailed       = errors.New("send failed")
	ErrNoData           = errors.New("no data")
	ErrNoDecrator       = errors.New("invalid process function for decorator")
	ErrInvalidMultiplex = errors.New("invalid process function for multiplexer")
	ErrInvalidSeparator = errors.New("invalid process function for separator")
	ErrUnknowType       = errors.New("unknown handler type")
)

type conveyerImpl struct {
	channels map[string]chan string
	handlers []handler
	mu       sync.Mutex
	size     int
}

type handler struct {
	process     interface{}
	handlerType string
	inputs      []string
	outputs     []string
}

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		channels: make(map[string]chan string),
		handlers: make([]handler, 0),
		mu:       sync.Mutex{},
		size:     size,
	}
}

func (conv *conveyerImpl) RegisterDecorator(
	handlerFunc func(ctx context.Context, input chan string, output chan string) error,
	input string, output string,
) {
	conv.mu.Lock()
	defer conv.mu.Unlock()
	conv.getOrCreateChannel(input)
	conv.getOrCreateChannel(output)
	conv.handlers = append(conv.handlers, handler{
		handlerType: "decorator",
		process:     handlerFunc,
		inputs:      []string{input},
		outputs:     []string{output},
	})
}

func (conv *conveyerImpl) RegisterMultiplexer(
	handlerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string, output string,
) {
	conv.mu.Lock()
	defer conv.mu.Unlock()

	for _, in := range inputs {
		conv.getOrCreateChannel(in)
	}

	conv.getOrCreateChannel(output)
	conv.handlers = append(conv.handlers, handler{
		handlerType: "multiplexer",
		process:     handlerFunc,
		inputs:      inputs,
		outputs:     []string{output},
	})
}

func (conv *conveyerImpl) RegisterSeparator(
	handlerFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	input string, outputs []string,
) {
	conv.mu.Lock()
	defer conv.mu.Unlock()
	conv.getOrCreateChannel(input)

	for _, out := range outputs {
		conv.getOrCreateChannel(out)
	}

	conv.handlers = append(conv.handlers, handler{
		handlerType: "separator",
		process:     handlerFunc,
		inputs:      []string{input},
		outputs:     outputs,
	})
}

func (conv *conveyerImpl) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var waitg sync.WaitGroup

	errorCh := make(chan error, 1)

	for _, handl := range conv.handlers {
		waitg.Add(1)

		go func(h handler) {
			defer waitg.Done()

			if err := conv.runHandler(ctx, h); err != nil {
				select {
				case errorCh <- err:
				default:
				}
				cancel()
			}
		}(handl)
	}

	select {
	case err := <-errorCh:
		if err != nil {
			return err
		}
	case <-ctx.Done():
		return fmt.Errorf("%w", ctx.Err())
	}

	return nil
}

func (conv *conveyerImpl) Send(input string, data string) error {
	conv.mu.Lock()
	chank, exists := conv.channels[input]
	conv.mu.Unlock()

	if !exists {
		return ErrChanNotFound
	}
	chank <- data

	return nil
}

func (conv *conveyerImpl) Recv(output string) (string, error) {
	conv.mu.Lock()
	channel, exists := conv.channels[output]
	conv.mu.Unlock()

	if !exists {
		return "", ErrChanNotFound
	}

	data, ok := <-channel
	if !ok {
		return "undefined", nil
	}

	return data, nil
}

func (conv *conveyerImpl) getOrCreateChannel(name string) chan string {
	if ch, exists := conv.channels[name]; exists {
		return ch
	}

	channel := make(chan string, conv.size)
	conv.channels[name] = channel

	return channel
}

func (conv *conveyerImpl) runHandler(ctx context.Context, hand handler) error {
	switch hand.handlerType {
	case "decorator":
		funct, ok := hand.process.(func(ctx context.Context, input chan string, output chan string) error)
		if !ok {
			return ErrNoDecrator
		}

		return funct(ctx, conv.channels[hand.inputs[0]], conv.channels[hand.outputs[0]])

	case "multiplexer":
		funct, ok := hand.process.(func(ctx context.Context, inputs []chan string, output chan string) error)
		if !ok {
			return ErrInvalidMultiplex
		}

		inputs := make([]chan string, len(hand.inputs))
		for i, in := range hand.inputs {
			inputs[i] = conv.channels[in]
		}

		return funct(ctx, inputs, conv.channels[hand.outputs[0]])

	case "separator":
		funct, ok := hand.process.(func(ctx context.Context, input chan string, outputs []chan string) error)
		if !ok {
			return ErrInvalidSeparator
		}

		outputs := make([]chan string, len(hand.outputs))
		for i, out := range hand.outputs {
			outputs[i] = conv.channels[out]
		}

		return funct(ctx, conv.channels[hand.inputs[0]], outputs)

	default:
		return ErrUnknowType
	}
}
