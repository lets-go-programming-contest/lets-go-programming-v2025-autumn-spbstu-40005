package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrNoDecorator  = errors.New("can't be decorated")
	ErrEmptyOutputs = errors.New("empty outputs")
	ErrEmptyInputs  = errors.New("empty inputs")
)

const (
	prefixValue      = "decorated: "
	skipDecoratorKey = "no decorator"
	skipMuxKey       = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case value, open := <-input:
			if !open {
				return nil
			}

			if strings.Contains(value, skipDecoratorKey) {
				return ErrNoDecorator
			}

			if !strings.HasPrefix(value, prefixValue) {
				value = prefixValue + value
			}

			select {
			case output <- value:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return ErrEmptyInputs
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(inputs))

	for _, inputCh := range inputs {
		ch := inputCh

		go func() {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case value, ok := <-ch:
					if !ok {
						return
					}

					if strings.Contains(value, skipMuxKey) {
						continue
					}

					select {
					case output <- value:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	waitGroup.Wait()

	return nil
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrEmptyOutputs
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case value, open := <-input:
			if !open {
				return nil
			}

			select {
			case outputs[index] <- value:
			case <-ctx.Done():
				return nil
			}

			index = (index + 1) % len(outputs)
		}
	}
}
