package handlers

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

var (
	ErrCannotBeDecorated = errors.New("can't be decorated")
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(data, "decorated: ") {
				data = "decorated: " + data
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- data:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, output := range outputs {
			close(output)
		}
	}()

	if len(outputs) == 0 {
		return nil
	}

	counter := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			outputIndex := counter % len(outputs)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case outputs[outputIndex] <- data:
			}

			counter++
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	readFromChannel := func(input chan string) {
		defer wg.Done()

		for {
			select {
			case <-subCtx.Done():
				return
			case data, ok := <-input:
				if !ok {
					return
				}

				if strings.Contains(data, "no multiplexer") {
					continue
				}

				select {
				case <-subCtx.Done():
					return
				case output <- data:
				}
			}
		}
	}

	for _, input := range inputs {
		wg.Add(1)
		go readFromChannel(input)
	}

	wg.Wait()
	return nil
}
