package handlers

import (
	"context"
	"errors"
	"strings"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case v, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(v, "no decorator") {
				return errors.New("can't be decorated")
			}

			if !strings.HasPrefix(v, "decorated: ") {
				v = "decorated: " + v
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- v:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {

	if len(outputs) == 0 {
		return nil
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case v, ok := <-input:
			if !ok {
				return nil
			}

			out := outputs[index]
			index = (index + 1) % len(outputs)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- v:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {

	if len(inputs) == 0 {
		return nil
	}

	openChannels := len(inputs)

	for openChannels > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			for i, ch := range inputs {
				select {
				case <-ctx.Done():
					return ctx.Err()

				case v, ok := <-ch:
					if !ok {
						inputs[i] = nil
						openChannels--
						continue
					}

					if strings.Contains(v, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return ctx.Err()
					case output <- v:
					}

				default:
				}
			}
		}
	}

	return nil
}
