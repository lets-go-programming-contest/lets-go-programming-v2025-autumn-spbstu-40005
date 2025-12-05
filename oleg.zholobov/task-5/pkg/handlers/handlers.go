package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

const (
	valDecoratedPrefix = "decorated: "
	valNoDecorator     = "no decorator"
	valNoMultiplexer   = "no multiplexer"
)

var (
	errCantDecorate = errors.New("can't be decorated")
	errEmptyOutputs = errors.New("outputs cannot be empty")
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {

	for {
		select {
		case <-ctx.Done():
			return nil

		case v, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(v, valNoDecorator) {
				return errCantDecorate
			}

			if !strings.HasPrefix(v, valDecoratedPrefix) {
				v = valDecoratedPrefix + v
			}

			select {
			case output <- v:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {

	if len(outputs) == 0 {
		return errEmptyOutputs
	}

	curIndex := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case v, ok := <-input:
			if !ok {
				return nil
			}

			idx := curIndex
			curIndex = (curIndex + 1) % len(outputs)

			select {
			case outputs[idx] <- v:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {

	if len(inputs) == 0 {
		return nil
	}

	wg := sync.WaitGroup{}
	wg.Add(len(inputs))

	for index := range inputs {
		go func(in chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case str, ok := <-in:
					if !ok {
						return
					}

					if strings.Contains(str, valNoMultiplexer) {
						continue
					}

					select {
					case output <- str:
					case <-ctx.Done():
						return
					}
				}
			}
		}(inputs[index])
	}

	wg.Wait()

	return nil
}
