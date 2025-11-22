package handlers

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
)

var (
	ErrCantBeDecorated = errors.New("can't be decorated")
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)

	prefix := "decorated: "

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(data, prefix) {
				data = prefix + data
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, output := range outputs {
			close(output)
		}
	}()

	if len(outputs) == 0 {
		return nil
	}

	var counter int64 = -1

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			index := atomic.AddInt64(&counter, 1) % int64(len(outputs))
			if index < 0 {
				index = 0
			}

			select {
			case outputs[index] <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	for {
		dataReceived := false
		for _, input := range inputs {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case data, ok := <-input:
				if !ok {
					continue
				}
				dataReceived = true
				
				if strings.Contains(data, "no multiplexer") {
					continue
				}

				select {
				case output <- data:
				case <-ctx.Done():
					return ctx.Err()
				}
			default:
			}
		}

		if !dataReceived {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
		}
	}
}
