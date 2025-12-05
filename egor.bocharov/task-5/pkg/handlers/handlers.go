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
	defer close(output) // Закрываем выходной канал при завершении

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil // Входной канал закрыт
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

	defer func() {
		// Закрываем все выходные каналы
		for _, out := range outputs {
			close(out)
		}
	}()

	index := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil // Входной канал закрыт
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
		close(output)
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	// Горутина для закрытия выходного канала после завершения всех обработчиков
	go func() {
		wg.Wait()
		close(output)
	}()

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
					case output <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	// Не ждем здесь, чтобы не блокировать
	return nil
}
