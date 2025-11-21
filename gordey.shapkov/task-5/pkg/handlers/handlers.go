package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrNoDecorator = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case str, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(str, "no decorator") {
				return ErrNoDecorator
			}

			if !strings.HasPrefix(str, "decorated: ") {
				str = "decorated: " + str
			}

			select {
			case output <- str:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wgroup sync.WaitGroup
	wgroup.Add(len(inputs))

	for _, ch := range inputs {
		go func(in chan string) {
			defer wgroup.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case str, ok := <-in:
					if !ok {
						return
					}
					if strings.Contains(str, "no multiplexer") {
						continue
					}
					select {
					case output <- str:
					case <-ctx.Done():
						return
					}
				}
			}
		}(ch)
	}

	wgroup.Wait()

	return nil
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		panic("empty outputs")
	}

	index := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case str, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outputs[index] <- str:
			case <-ctx.Done():
				return nil
			}
			index = (index + 1) % len(outputs)
		}
	}
}
