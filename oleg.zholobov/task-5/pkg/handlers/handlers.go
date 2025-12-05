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
		case value, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(value, valNoDecorator) {
				return errCantDecorate
			}

			if !strings.HasPrefix(value, valDecoratedPrefix) {
				value = valDecoratedPrefix + value
			}

			select {
			case output <- value:
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
		case value, ok := <-input:
			if !ok {
				return nil
			}

			idx := curIndex
			curIndex = (curIndex + 1) % len(outputs)

			select {
			case outputs[idx] <- value:
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

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(inputs))

	for index := range inputs {
		go func(inputChannel chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case str, ok := <-inputChannel:
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

	waitGroup.Wait()

	return nil
}
