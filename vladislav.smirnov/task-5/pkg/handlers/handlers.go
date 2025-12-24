package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrCantDecorate = errors.New("can't be decorated")
	ErrEmptyChannel = errors.New("empty channel")
)

const (
	prefix        = "decorated: "
	noDecorator   = "no decorator"
	noMultiplexer = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case value, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(value, noDecorator) {
				return ErrCantDecorate
			}

			if !strings.HasPrefix(value, prefix) {
				value = prefix + value
			}

			select {
			case output <- value:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrEmptyChannel
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case value, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outputs[index] <- value:
				index = (index + 1) % len(outputs)

			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return ErrEmptyChannel
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputs))

	for _, input := range inputs {
		go func(in <-chan string) {
			defer waitGroup.Done()
			processChannel(ctx, in, output)
		}(input)
	}

	waitGroup.Wait()

	return nil
}

func processChannel(ctx context.Context, input <-chan string, output chan<- string) {
	for {
		select {
		case <-ctx.Done():
			return

		case value, ok := <-input:
			if !ok {
				return
			}

			if strings.Contains(value, noMultiplexer) {
				continue
			}

			select {
			case output <- value:

			case <-ctx.Done():
			}
		}
	}
}
