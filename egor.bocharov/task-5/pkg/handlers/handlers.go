package handlers

import (
	"context"
	"errors"
	"strings"
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

			if strings.Contains(data, noDecorator) {
				return ErrDecorator
			}

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
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			allClosed := true

			for _, inputCh := range inputs {
				select {
				case data, ok := <-inputCh:
					if !ok {
						continue
					}

					allClosed = false

					if strings.Contains(data, noMultiplexer) {
						continue
					}

					select {
					case output <- data:
					case <-ctx.Done():
						return nil
					}
				default:
					// continue
				}
			}

			if allClosed {
				return nil
			}
		}
	}
}
