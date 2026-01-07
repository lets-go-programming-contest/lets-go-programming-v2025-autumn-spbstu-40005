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
		case msg, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(msg, noDecoratorMark) {
				return ErrCantBeDecorated
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
		case msg, ok := <-input:
			if !ok {
				return nil
			}

			outCh := outputs[index]
			index++
			if index >= len(outputs) {
				index = 0
			}

			select {
			case <-ctx.Done():
				return nil
			case outCh <- msg:
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

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, in := range inputs {
		inCh := in
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-inCh:
					if !ok {
						return
					}

					if strings.Contains(msg, noMultiplexMark) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- msg:
					}
				}
			}
		}()
	}

	wg.Wait()
	return nil
}
