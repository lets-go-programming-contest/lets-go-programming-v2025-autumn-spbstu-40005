package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrCantDecorate     = errors.New("can't be decorated")
	ErrNoOutputChannels = errors.New("no output channels")
	ErrNoInputChannels  = errors.New("no input channels")
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	const prefix = "decorated: "

	const noDecorator = "no decorator"

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, noDecorator) {
				return ErrCantDecorate
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
	if len(outputs) == 0 {
		return ErrNoOutputChannels
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case data, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outputs[index] <- data:
				index = (index + 1) % len(outputs)

			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	const skipMultiplexer = "no multiplexer"

	if len(inputs) == 0 {
		return ErrNoInputChannels
	}

	var wg sync.WaitGroup

	wg.Add(len(inputs))

	for _, in := range inputs {
		go func(inputChan chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case data, ok := <-inputChan:
					if !ok {
						return
					}

					if strings.Contains(data, skipMultiplexer) {
						continue
					}

					select {
					case output <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(in)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
