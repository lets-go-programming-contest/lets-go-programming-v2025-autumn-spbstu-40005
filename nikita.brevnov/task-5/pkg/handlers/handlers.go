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

func AddPrefix(ctx context.Context, in chan string, out chan string) error {
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

func Distribute(ctx context.Context, input chan string, outs []chan string) error {
	if len(outs) == 0 {
		return ErrNoOutputs
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, isOpen := <-input:
			if !isOpen {
				return nil
			}

			outs[index%len(outs)] <- data
			index++
		}
	}
}

func Merge(ctx context.Context, ins []chan string, out chan string) error {
	var waitGroup sync.WaitGroup

	for _, inputChan := range ins {
		waitGroup.Add(1)

		go func(inputCh chan string) {
			defer waitGroup.Done()

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
		}(inputChan)
	}

	waitGroup.Wait()

	return nil
}
