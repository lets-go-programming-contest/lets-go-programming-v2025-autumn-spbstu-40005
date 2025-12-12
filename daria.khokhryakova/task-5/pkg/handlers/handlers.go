package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrProcessingFailed = errors.New("can't be decorated")
	ErrInvalidConfig    = errors.New("invalid handler configuration")
)

const (
	prefix        = "decorated: "
	noDecorator   = "no decorator"
	noMultiplexer = "no multiplexer"
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

			if strings.Contains(data, noDecorator) {
				return ErrProcessingFailed
			}

			if !strings.HasPrefix(data, prefix) {
				data = prefix + data
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- data:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrInvalidConfig
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			outputChan := outputs[index%len(outputs)]
			index++

			select {
			case <-ctx.Done():
				return nil
			case outputChan <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return ErrInvalidConfig
	}

	var waitgr sync.WaitGroup

	waitgr.Add(len(inputs))

	for _, input := range inputs {
		go func(channel chan string) {
			defer waitgr.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-channel:
					if !ok {
						return
					}

					if strings.Contains(data, noMultiplexer) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- data:
					}
				}
			}
		}(input)
	}

	waitgr.Wait()

	return nil
}
