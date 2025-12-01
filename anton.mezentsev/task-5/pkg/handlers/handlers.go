package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrProcessingFailed = errors.New("processing failed")
	ErrNoOutputs        = errors.New("outputs cannot be empty")
)

const (
	skipDecorator   = "no decorator"
	decoration      = "decorated: "
	skipMultiplexer = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, in chan string, out chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, isOpen := <-in:
			if !isOpen {
				return nil
			}

			if strings.Contains(data, skipDecorator) {
				return ErrProcessingFailed
			}

			if !strings.HasPrefix(data, decoration) {
				data = decoration + data
			}

			out <- data
		}
	}
}

func SeparatorFunc(ctx context.Context, in chan string, outs []chan string) error {
	if len(outs) == 0 {
		return ErrNoOutputs
	}

	roundRobinIndex := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, isOpen := <-in:
			if !isOpen {
				return nil
			}

			targetCh := outs[roundRobinIndex%len(outs)]
			roundRobinIndex++

			targetCh <- data
		}
	}
}

func MultiplexerFunc(ctx context.Context, ins []chan string, out chan string) error {
	var wg sync.WaitGroup

	for _, ch := range ins {
		wg.Add(1)

		go func(inputCh chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, isOpen := <-inputCh:
					if !isOpen {
						return
					}

					if strings.Contains(data, skipMultiplexer) {
						continue
					}

					out <- data
				}
			}
		}(ch)
	}

	wg.Wait()
	return nil
}
