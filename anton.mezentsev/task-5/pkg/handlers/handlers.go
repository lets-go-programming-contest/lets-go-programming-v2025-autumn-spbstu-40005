package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrDecorator          = errors.New("can't be decorated")
	ErrNoOutputs          = errors.New("outputs cannot be empty")
	ErrNoInputs           = errors.New("inputs cannot be empty")
	ErrInputChannelClosed = errors.New("input channel closed")
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
				return ErrInputChannelClosed
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

	roundRobinIndex := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return ErrInputChannelClosed
			}

			targetChannel := outputs[roundRobinIndex%len(outputs)]
			roundRobinIndex++

			select {
			case targetChannel <- data:
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

	var waitGroup sync.WaitGroup
	errCh := make(chan error, 1)

	for _, inputChannel := range inputs {
		waitGroup.Add(1)

		go func(channel chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-channel:
					if !ok {
						select {
						case errCh <- ErrInputChannelClosed:
						default:
						}
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
		}(inputChannel)
	}

	go func() {
		waitGroup.Wait()
		close(errCh)
	}()

	if err := <-errCh; err != nil {
		return err
	}

	return nil
}
