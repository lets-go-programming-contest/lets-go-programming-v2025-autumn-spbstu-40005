package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrDecorator     = errors.New("can't be decorated")
	ErrOutputsEmpty  = errors.New("outputs must not be empty")
	ErrInputsEmpty   = errors.New("inputs must not be empty")
)

const (
	noDecorator   = "no decorator"
	prefix        = "decorated: "
	noMultiplexer = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, noDecorator) {
				return ErrDecorator
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
		return ErrOutputsEmpty
	}

	counter := 0
	outputsCount := len(outputs)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			idx := counter % outputsCount
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
		return ErrInputsEmpty
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	for _, inputCh := range inputs {
		wg.Add(1)

		go func(ch chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-ch:
					if !ok {
						return
					}

					if strings.Contains(data, noMultiplexer) {
						continue
					}

					select {
					case output <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(inputCh)
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
	case err := <-errChan:
		return err
	}
}
