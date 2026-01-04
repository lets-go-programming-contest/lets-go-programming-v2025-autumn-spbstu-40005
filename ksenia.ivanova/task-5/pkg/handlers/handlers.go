package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	const prefix = "decorated: "
	const errorSubstring = "no decorator"
	const errorMsg = "can't be decorated"

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, errorSubstring) {
				return errors.New(errorMsg)
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
		return nil
	}

	counter := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
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
	const skipSubstring = "no multiplexer"

	if len(inputs) == 0 {
		return nil
	}

	done := ctx.Done()
	errChan := make(chan error, len(inputs))
	var wg sync.WaitGroup

	for _, in := range inputs {
		wg.Add(1)
		go func(inputChan chan string) {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case data, ok := <-inputChan:
					if !ok {
						return
					}
					if strings.Contains(data, skipSubstring) {
						continue
					}
					select {
					case output <- data:
					case <-done:
						return
					}
				}
			}
		}(in)
	}

	go func() {
		wg.Wait()
		select {
		case errChan <- nil:
		default:
		}
	}()

	select {
	case <-done:
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}
