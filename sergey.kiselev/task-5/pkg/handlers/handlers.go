package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrDecorator    = errors.New("can't be decorated")
	ErrOutputsEmpty = errors.New("outputs must not be empty")
	ErrInputsEmpty  = errors.New("inputs must not be empty")
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
			return nil
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
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrOutputsEmpty
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

			outCh := outputs[index%len(outputs)]
			index++

			select {
			case outCh <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return ErrInputsEmpty
	}
	var waitgrr sync.WaitGroup

	waitgrr.Add(len(inputs))

	for _, inputCh := range inputs {

		go func(channel chan string) {
			defer waitgrr.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-channel:
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

	waitgrr.Wait()

	return nil
}
