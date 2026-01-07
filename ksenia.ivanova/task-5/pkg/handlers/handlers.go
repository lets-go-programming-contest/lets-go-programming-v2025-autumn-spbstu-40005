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
	const (
		prefix      = "decorated: "
		noDecorator = "no decorator"
	)

	for {
		select {
		case <-ctx.Done():
			return nil
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
				return nil
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
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			channel := outputs[index]
			index = (index + 1) % len(outputs)

			select {
			case channel <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	const skipMultiplexer = "no multiplexer"

	if len(inputs) == 0 {
		return ErrNoInputChannels
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputs))

	for _, input := range inputs {
		go func(channel chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case data, ok := <-channel:
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
		}(input)
	}

	waitGroup.Wait()

	return nil
}
