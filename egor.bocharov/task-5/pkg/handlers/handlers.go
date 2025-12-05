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

// PrefixDecoratorFunc добавляет префикс к данным, если его ещё нет.
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

			outChannel := outputs[index%len(outputs)]
			index++

			select {
			case outChannel <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

// MultiplexerFunc объединяет данные из нескольких входных каналов в один выходной.
// Данные с пометкой "no multiplexer" пропускаются.
func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, ch := range inputs {
		inputCh := ch // захват для горутины
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-inputCh:
					if !ok {
						return // канал закрыт
					}
					if strings.Contains(data, noMultiplexer) {
						continue
					}
					select {
					case output <- data: // ✅ ИСПРАВЛЕНО: добавлено 'data'
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	wg.Wait()
	return nil
}
