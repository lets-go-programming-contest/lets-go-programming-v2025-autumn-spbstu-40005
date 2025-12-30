package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

const (
	decoratedPrefix = "decorated: "
	stopDecorator   = "no decorator"
	stopMultiplexer = "no multiplexer"
)

var (
	ErrCannotBeDecorated = errors.New("can't be decorated")
	ErrEmptyOutputs      = errors.New("empty outputs")
	ErrEmptyInputs       = errors.New("empty inputs")
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(msg, stopDecorator) {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(msg, decoratedPrefix) {
				msg = decoratedPrefix + msg
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- msg:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return ErrEmptyInputs
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputs))

	for _, ch := range inputs {
		inch := ch
		go func(inputChan chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-inputChan:
					if !ok {
						return
					}

					if strings.Contains(msg, stopMultiplexer) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- msg:
					}
				}
			}
		}(inch)
	}

	waitGroup.Wait()

	return nil
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrEmptyOutputs
	}

	idx := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-input:
			if !ok {
				return nil
			}

			target := outputs[idx]
			idx = (idx + 1) % len(outputs)
			select {
			case <-ctx.Done():
				return nil
			case target <- msg:
			}
		}
	}
}
