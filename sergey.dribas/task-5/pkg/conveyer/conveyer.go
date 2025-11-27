package conveyer

import (
	"context"
	"errors"
	"sync"
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
		size:     size,
	}
}

func (conv *conveyerImpl) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string, output string,
) {
	conv.mu.Lock()
	defer conv.mu.Unlock()
	conv.getOrCreateChannel(input)
	conv.getOrCreateChannel(output)
	conv.handlers = append(conv.handlers, handler{
		handlerType: "decorator",
		process:     fn,
		inputs:      []string{input},
		outputs:     []string{output},
	})
}

func (conv *conveyerImpl) RegisterMultipleser(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
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
		process:     fn,
		inputs:      inputs,
		outputs:     []string{output},
	})
}

func (conv *conveyerImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
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
		process:     fn,
		inputs:      []string{input},
		outputs:     outputs,
	})
}

func (conv *conveyerImpl) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errorCh := make(chan error, len(conv.handlers))

	for _, h := range conv.handlers {
		wg.Add(1)
		go func(h handler) {
			defer wg.Done()
			if err := conv.runHandler(ctx, h); err != nil {
				select {
				case errorCh <- err:
				default:
				}
				cancel()
			}
		}(h)
	}

	wg.Wait()
	close(errorCh)

	for err := range errorCh {
		if err != nil {
			return err
		}
	}
	return nil
}

func (conv *conveyerImpl) Send(input string, data string) error {
	conv.mu.Lock()
	ch, exists := conv.channels[input]
	conv.mu.Unlock()
	if !exists {
		return errors.New("chan not found")
	}
	select {
	case ch <- data:
		return nil
	default:
		return errors.New("send failed")
	}
}

func (conv *conveyerImpl) Recv(output string) (string, error) {
	conv.mu.Lock()
	ch, exists := conv.channels[output]
	conv.mu.Unlock()
	if !exists {
		return "", errors.New("chan not found")
	}
	select {
	case data, ok := <-ch:
		if !ok {
			return "undefined", nil
		}
		return data, nil
	default:
		return "", errors.New("no data")
	}
}

func (conv *conveyerImpl) getOrCreateChannel(name string) chan string {
	if ch, exists := conv.channels[name]; exists {
		return ch
	}
	ch := make(chan string, conv.size)
	conv.channels[name] = ch
	return ch
}

func (conv *conveyerImpl) runHandler(ctx context.Context, h handler) error {
	switch h.handlerType {
	case "decorator":
		fn := h.process.(func(ctx context.Context, input chan string, output chan string) error)
		return fn(ctx, conv.channels[h.inputs[0]], conv.channels[h.outputs[0]])
	case "multiplexer":
		fn := h.process.(func(ctx context.Context, inputs []chan string, output chan string) error)
		inputs := make([]chan string, len(h.inputs))
		for i, in := range h.inputs {
			inputs[i] = conv.channels[in]
		}
		return fn(ctx, inputs, conv.channels[h.outputs[0]])
	case "separator":
		fn := h.process.(func(ctx context.Context, input chan string, outputs []chan string) error)
		outputs := make([]chan string, len(h.outputs))
		for i, out := range h.outputs {
			outputs[i] = conv.channels[out]
		}
		return fn(ctx, conv.channels[h.inputs[0]], outputs)
	}
	return nil
}
