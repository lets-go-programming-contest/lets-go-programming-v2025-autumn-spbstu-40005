package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrNoDecorator = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for str := range input {
		if strings.Contains(str, "no decorator") {
			return ErrNoDecorator
		}

		if !strings.HasPrefix(str, "decorated: ") {
			str = "decorated: " + str
		}
		select {
		case <-ctx.Done():
			return nil

		case output <- str:
		}
	}

	return nil
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wgroup sync.WaitGroup
	wgroup.Add(len(inputs))

	for _, ch := range inputs {
		go func(channel chan string) {
			defer wgroup.Done()

			for str := range channel {
				if !strings.Contains(str, "no multiplexer") {
					select {
					case <-ctx.Done():
						return

					case output <- str:
					}
				} else {
					continue
				}
			}
		}(ch)

	}

	wgroup.Wait()

	return nil
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	length := len(outputs)
	if length == 0 {
		panic("empty outputs")
	}

	index := 0
	for str := range input {
		select {
		case <-ctx.Done():
			return nil

		case outputs[index] <- str:
		}

		index = (index + 1) % length
	}

	return nil
}
