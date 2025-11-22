package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context done: %w", ctx.Err())
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(data, "decorated: ") {
				data = "decorated: " + data
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- data:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
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
		return nil
	}

	waitGroup := sync.WaitGroup{}
	subContext, cancel := context.WithCancel(ctx)

	defer cancel()

	readFromChannel := func(inputChannel chan string) {
		defer waitGroup.Done()

		for {
			select {
			case <-subContext.Done():
				return
			case data, ok := <-inputChannel:
				if !ok {
					return
				}

				if strings.Contains(data, "no multiplexer") {
					continue
				}

				select {
				case <-subContext.Done():
					return
				case output <- data:
				}
			}
		}
	}

	for _, inputChannel := range inputs {
		waitGroup.Add(1)

		channel := inputChannel
		go readFromChannel(channel)
	}

	waitGroup.Wait()

	return nil
}
