package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

const (
	decorator_prefix    = "decorated: "
	no_decorator_prefix = "no decorator"
	no_multip_str       = "no multiplexer"
)

var ErrCantDecorate = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, no_decorator_prefix) {
				return ErrCantDecorate
			}

			if !strings.HasPrefix(data, decorator_prefix) {
				data = decorator_prefix + data
			}

			output <- data
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	counter := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			idx := counter % len(outputs)
			counter++

			select {
			case outputs[idx] <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	var wg sync.WaitGroup

	for _, input := range inputs {
		wg.Add(1)
		go func(in chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-in:
					if !ok {
						return
					}

					if !strings.Contains(data, no_multip_str) {
						select {
						case output <- data:
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}(input)
	}

	wg.Wait()
	return nil
}
