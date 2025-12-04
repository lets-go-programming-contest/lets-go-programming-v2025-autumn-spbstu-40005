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

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	const prefix = "decorated: "

	for {
		select {
		case <-ctx.Done():
			return nil

		case value, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(value, "no decorator") {
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

	for _, channel := range inputs {
		go func(in <-chan string) {
			defer waitGroup.Done()
			processString(ctx, in, output)
		}(input)
	}

	waitGroup.Wait()

	return nil
}

func processChannel(ctx context.Context, input <-chan string, output chan<- string) {
	for {
		select {
		case <-gCtx.Done():
			return gCtx.Err()

		case value, ok := <-channel:
			if !ok {
				return nil
			}

			if strings.Contains(value, "no nultiplexer") {
				continue
			}

			select {
			case output <- value:

			case <-gCtx.Done():
				return gCtx.Err()
			}
		}
	}
}
