package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrDecorator    = errors.New("can't be decorated")
	ErrOutputsEmpty = errors.New("outputs must not be empty")
)

const (
	noDecorator   = "no decorator"
	prefix        = "decorated: "
	noMultiplexer = "no multiplexer"
)

// PrefixDecoratorFunc добавляет префикс к данным, если его еще нет.
// Если данные содержат "no decorator", возвращает ошибку.
func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			containsNoDecorator := strings.Contains(data, noDecorator)
			if containsNoDecorator {
				return ErrDecorator
			}

			hasPrefix := strings.HasPrefix(data, prefix)
			if !hasPrefix {
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

// SeparatorFunc распределяет данные по выходным каналам по кругу.
func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrOutputsEmpty
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

			outCh := outputs[index%len(outputs)]
			index++

			select {
			case outCh <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

// MultiplexerFunc объединяет данные из нескольких входных каналов в один выходной.
// Данные с пометкой "no multiplexer" пропускаются.
func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var waitGroup sync.WaitGroup

	workerFunc := func(inputChannel chan string) {
		defer waitGroup.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-inputChannel:
				if !ok {
					return
				}

				containsNoMultiplexer := strings.Contains(data, noMultiplexer)
				if containsNoMultiplexer {
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

	for _, inputCh := range inputs {
		waitGroup.Add(1)
		worker := workerFunc
		go worker(inputCh)
	}

	waitGroup.Wait()

	return nil
}
