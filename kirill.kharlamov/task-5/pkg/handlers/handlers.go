package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrDecorator    = errors.New("data cannot be decorated")
	ErrOutputsEmpty = errors.New("outputs slice must not be empty")
)

const (
	noDecorator   = "no decorator"
	prefix        = "decorated: "
	noMultiplexer = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, inputChannel, outputChannel chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-inputChannel:
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
			case outputChannel <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, inputChannel chan string, outputChannels []chan string) error {
	if len(outputChannels) == 0 {
		return ErrOutputsEmpty
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-inputChannel:
			if !ok {
				return nil
			}

			outChannel := outputChannels[index%len(outputChannels)]
			index++

			select {
			case outChannel <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputChannels []chan string, outputChannel chan string) error {
	var waitGroup sync.WaitGroup

	worker := func(channel chan string) {
		defer waitGroup.Done()

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
				case outputChannel <- data:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	for _, inputChan := range inputChannels {
		waitGroup.Add(1)
		go worker(inputChan)
	}

	waitGroup.Wait()

	return nil
}
