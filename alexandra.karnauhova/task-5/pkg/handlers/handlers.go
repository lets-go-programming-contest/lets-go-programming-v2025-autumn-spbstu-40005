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

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return ErrCantDecorate
			}

			prefix := "decorated: "
			if !strings.HasPrefix(data, prefix) {
				data = prefix + data
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var waitg sync.WaitGroup

	reader := func(channel chan string) {
		defer waitg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-channel:
				if !ok {
					return
				}

				if strings.Contains(data, "no multiplexer") {
					continue
				}

				select {
				case output <- data:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	for _, channel := range inputs {
		waitg.Add(1)

		go reader(channel)
	}

	waitg.Wait()

	return nil
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	idx := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			target := outputs[idx%len(outputs)]
			idx++

			select {
			case target <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
