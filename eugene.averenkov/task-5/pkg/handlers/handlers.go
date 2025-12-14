package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrCannotBeDecorated = errors.New("can't be decorated")
	ErrInput = errors.New("no input channels provided")
	ErrOutput = errors.New("no output channels provided")
)

const (
	noDecoratorPrefix   = "no decorator"
	decoratedPrefix     = "decorated: "
	noMultiplexerString = "no multiplexer"
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

			if strings.Contains(data, noDecoratorPrefix) {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(data, decoratedPrefix) {
				data = decoratedPrefix + data
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
		return ErrOutput
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

			outputIndex := counter % len(outputs)

			select {
			case <-ctx.Done():
				return nil
			case outputs[outputIndex] <- data:
			}

			counter++
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return ErrInput
	}

	waitGroup := sync.WaitGroup{}

	for _, inputChannel := range inputs {
		waitGroup.Add(1)

		channel := inputChannel
		readFunc := func() {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-channel:
					if !ok {
						return
					}

					if strings.Contains(data, noMultiplexerString) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- data:
					}
				}
			}
		}

		go readFunc()
	}

	waitGroup.Wait()
	return nil
}
