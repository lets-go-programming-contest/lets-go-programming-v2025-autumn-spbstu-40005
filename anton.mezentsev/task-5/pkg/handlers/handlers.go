package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrDecorator = errors.New("can't be decorated")
	ErrNoOutputs = errors.New("outputs cannot be empty")
	ErrNoInputs  = errors.New("inputs cannot be empty")
)

const (
	skipDecorator   = "no decorator"
	decoration      = "decorated: "
	skipMultiplexer = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {

				return nil
			}

			if strings.Contains(data, skipDecorator) {
				return ErrDecorator
			}

			if !strings.HasPrefix(data, decoration) {
				data = decoration + data
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrNoOutputs
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {

				return nil
			}

			target := outputs[index]
			index = (index + 1) % len(outputs)

			select {
			case target <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return ErrNoInputs
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	for _, in := range inputs {
		wg.Add(1)

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

	go func() {
		wg.Wait()
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return nil
	}
}
