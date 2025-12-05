package handlers

import (
	"context"
	"errors"
	"fmt"
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
			return fmt.Errorf("prefix decorator: %w", ctx.Err())
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
				return fmt.Errorf("prefix decorator: %w", ctx.Err())
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
			return fmt.Errorf("separator: %w", ctx.Err())
		case data, ok := <-inputChannel:
			if !ok {
				return nil
			}

			selectedOutputChannel := outputChannels[index%len(outputChannels)]
			index++

			select {
			case selectedOutputChannel <- data:
			case <-ctx.Done():
				return fmt.Errorf("separator: %w", ctx.Err())
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputChannels []chan string, outputChannel chan string) error {
	var waitGroup sync.WaitGroup

	worker := func(inputChannel chan string) {
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
	}

	for _, inputChan := range inputChannels {
		waitGroup.Add(1)

		currentInputChan := inputChan
		workerFunc := worker

		go func() {
			workerFunc(currentInputChan)
		}()
	}

	waitGroup.Wait()

	select {
	case <-ctx.Done():
		return fmt.Errorf("multiplexer: %w", ctx.Err())
	default:
		return nil
	}
}
