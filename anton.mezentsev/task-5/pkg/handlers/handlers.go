package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

const (
	tagDecorated    = "decorated: "
	stopDecorator   = "no decorator"
	stopMultiplexer = "no multiplexer"
)

var (
	ErrDecorationRefused  = errors.New("can't be decorated")
	ErrEmptyOutputs       = errors.New("empty outputs channels")
	ErrEmptyInputs        = errors.New("empty inputs channels")
	ErrInputChannelClosed = errors.New("input channel closed")
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

			if strings.Contains(message, stopDecorator) {
				return ErrDecorationRefused
			}

			if !strings.HasPrefix(message, tagDecorated) {
				message = tagDecorated + message
			}

			select {
			case output <- message:
			case <-ctx.Done():
				return nil
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
		return ErrEmptyInputs
	}

	var waitGrp sync.WaitGroup

	waitGrp.Add(len(inputs))

	for _, channel := range inputs {
		go func(inputChannel chan string) {
			defer waitGrp.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case message, ok := <-inputChannel:
					if !ok {
						return
					}

					if strings.Contains(message, stopMultiplexer) {
						continue
					}

					select {
					case output <- message:
					case <-ctx.Done():
						return
					}
				}
			}
		}(channel)
	}

	waitGrp.Wait()

	return nil
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	if len(outputs) == 0 {
		return ErrEmptyOutputs
	}

	var assignmentIndex int

	for {
		select {
		case <-ctx.Done():
			return nil
		case message, ok := <-input:
			if !ok {
				return nil
			}

			targetChannel := outputs[assignmentIndex]
			assignmentIndex = (assignmentIndex + 1) % len(outputs)

			select {
			case targetChannel <- message:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
