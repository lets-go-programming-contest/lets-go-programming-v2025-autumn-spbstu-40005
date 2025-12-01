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
		case data, isOpen := <-input:
			if !isOpen {
				return nil
			}

			if strings.Contains(data, skipDecorator) {
				return ErrDecorator
			}

			if !strings.HasPrefix(data, decoration) {
				data = decoration + data
			}

			output <- data
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
		case data, isOpen := <-input:
			if !isOpen {
				return nil
			}

			targetChannel := outputs[roundRobinIndex%len(outputs)]
			roundRobinIndex++

			targetChannel <- data
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var waitGroup sync.WaitGroup

	for _, inputChannel := range inputs {
		waitGroup.Add(1)

		go func(channel chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, isOpen := <-channel:
					if !isOpen {
						return
					}

					if strings.Contains(data, skipMultiplexer) {
						continue
					}

					output <- data
				}
			}
		}(inputChannel)
	}

	waitGroup.Wait()

	return nil
}
