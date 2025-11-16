package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var (
	ErrNoDecorator = errors.New("can’t be decorated")
)

func orDone[T any](done <-chan struct{}, channel <-chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return

			default:
			}

			select {
			case <-done:
				return

			case v, ok := <-channel:
				if !ok {
					return
				}

				select {
				case out <- v:

				case <-done:
					return
				}
			}
		}
	}()

	return out
}

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	for str := range orDone(ctx.Done(), input) {
		if strings.Contains(str, "no decorator") {
			return ErrNoDecorator
		}

		if !strings.HasPrefix(str, "decorated: ") {
			str = fmt.Sprintf("decorated: %s", str)
		}

		select {
		case output <- str:

		case <-ctx.Done():
			return nil
		}
	}

	return nil
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	index := 0

	for str := range orDone(ctx.Done(), input) {
		index = (index + 1) % len(outputs)

		select {
		case outputs[index] <- str:

		case <-ctx.Done():
			return nil
		}
	}

	return nil
}

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	wg := sync.WaitGroup{}

	wg.Add(len(inputs))

	for _, ch := range inputs {
		go func() {
			defer wg.Done()

			for str := range orDone(ctx.Done(), ch) {
				if strings.Contains(str, "no multiplexer") {
					continue
				}

				select {
				case output <- str:

				case <-ctx.Done():
					return
				}
			}
		}()
	}

	wg.Wait()

	return nil
}
