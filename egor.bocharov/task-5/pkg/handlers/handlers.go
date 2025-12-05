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

// PrefixDecoratorFunc добавляет префикс к данным, если его ещё нет.
func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer func() {
		// Пытаемся закрыть выходной канал
		select {
		case <-output:
		default:
			close(output)
		}
	}()

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

	defer func() {
		// Закрываем все выходные каналы
		for _, out := range outputs {
			select {
			case <-out:
			default:
				close(out)
			}
		}
	}()

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
func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		select {
		case <-output:
		default:
			close(output)
		}
		return nil
	}

	// Создаем отдельный контекст для этой функции
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Создаем канал для данных
	dataChan := make(chan string, len(inputs))

	// Запускаем горутины для чтения из всех входных каналов
	for _, in := range inputs {
		go func(inputChan chan string) {
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-inputChan:
					if !ok {
						return
					}
					if strings.Contains(data, noMultiplexer) {
						continue
					}
					select {
					case dataChan <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(in)
	}

	// Главная горутина записывает в выходной канал
	defer func() {
		select {
		case <-output:
		default:
			close(output)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-dataChan:
			if !ok {
				return nil
			}
			select {
			case output <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
