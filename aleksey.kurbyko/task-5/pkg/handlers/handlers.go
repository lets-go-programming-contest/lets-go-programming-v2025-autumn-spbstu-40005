package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrCantBeDecorated = errors.New("can't be decorated")
	ErrEmptyOutputs    = errors.New("empty outputs")
)

const (
	decoratedPrefix = "decorated: "
	noDecoratorMark = "no decorator"
	noMultiplexMark = "no multiplexer"
)

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case message, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(message, noDecoratorMark) {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(message, decoratedPrefix) {
				message = decoratedPrefix + message
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- message:
			}
		}
	}
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	if len(outputs) == 0 {
		return ErrEmptyOutputs
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case message, ok := <-input:
			if !ok {
				return nil
			}

			outputChannel := outputs[index]
			index = (index + 1) % len(outputs)

			select {
			case <-ctx.Done():
				return nil
			case outputChannel <- message:
			}
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	if len(inputs) == 0 {
		return nil
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputs))

	for _, inputChannel := range inputs {
		currentInput := inputChannel

		worker := func() {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case message, ok := <-currentInput:
					if !ok {
						return
					}

					if strings.Contains(message, noMultiplexMark) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- message:
					}
				}
			}
		}

		go worker()
	}

	waitGroup.Wait()

	return nil
}
