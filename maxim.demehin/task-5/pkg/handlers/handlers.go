package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

const (
	decoratorPrefix   = "decorated: "
	noDecoratorPrefix = "no decorator"
	noMultipStr       = "no multiplexer"
)

var (
	ErrCantDecorate = errors.New("can't be decorated")
	ErrEmptyChannel = errors.New("empty channel")
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, noDecoratorPrefix) {
				return ErrCantDecorate
			}

			if !strings.HasPrefix(data, decoratorPrefix) {
				data = decoratorPrefix + data
			}

			output <- data
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	outsLen := len(outputs)

	if outsLen == 0 {
		return ErrEmptyChannel
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

			idx := counter % outsLen
			counter++

			select {
			case outputs[idx] <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	inputsLen := len(inputs)
	if inputsLen == 0 {
		return ErrEmptyChannel
	}

	var waitGr sync.WaitGroup

	waitGr.Add(inputsLen)

	for _, input := range inputs {
		go processInputChannel(ctx, &waitGr, input, output)
	}

	waitGr.Wait()

	return nil
}

func processInputChannel(ctx context.Context, waitGr *sync.WaitGroup, input <-chan string, output chan<- string) {
	defer waitGr.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-input:
			if !ok {
				return
			}

			if strings.Contains(data, noMultipStr) {
				continue
			}

			select {
			case output <- data:
			case <-ctx.Done():
			}
		}
	}
}
