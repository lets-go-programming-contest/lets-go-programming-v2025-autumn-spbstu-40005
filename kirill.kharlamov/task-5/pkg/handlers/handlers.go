package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrDecorator    = errors.New("can't be decorated")
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

			selectedOutputChannel := outputChannels[index%len(outputChannels)]
			index++

			select {
			case selectedOutputChannel <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputChannels []chan string, outputChannel chan string) error {
	if len(inputChannels) == 0 {
		return nil
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(inputChannels))

	for _, inputChan := range inputChannels {
		go func(inputChannel chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-inputChannel:
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
		}(inputChan)
	}

	waitGroup.Wait()

	return nil
}
