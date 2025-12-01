package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
)

var ErrProcessingFailed = errors.New("can't be decorated")

const (
	prefix        = "decorated: "
	noDecorator   = "no decorator"
	noMultiplexer = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, noDecorator) {
				return ErrProcessingFailed
			}

			if !strings.HasPrefix(data, prefix) {
				data = prefix + data
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
	if len(outputs) == 0 {
		return nil
	}

	var counter uint64

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			index := atomic.AddUint64(&counter, 1) % uint64(len(outputs))

			select {
			case <-ctx.Done():
				return ctx.Err()
			case outputs[index] <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	merged := make(chan string)

	for _, input := range inputs {
		wg.Add(1)
		go func(in chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-in:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case merged <- data:
					}
				}
			}
		}(input)
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-merged:
			if !ok {
				return nil
			}
			if strings.Contains(data, noMultiplexer) {
				continue
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- data:
			}
		}
	}
}
