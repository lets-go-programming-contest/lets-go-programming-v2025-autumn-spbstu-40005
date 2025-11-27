package handlers

import (
	"context"
	"errors"
	"strings"
	"time"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case data, ok := <-input:
			if !ok {
				return nil
			}
			if strings.Contains(data, "no decorator") {
				return errors.New("can't be decorated")
			}
			if !strings.HasPrefix(data, "decorated: ") {
				data = "decorated: " + data
			}
			select {
			case output <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	counter := 0
	for {
		select {
		case data, ok := <-input:
			if !ok {
				return nil
			}
			idx := counter % len(outputs)
			select {
			case outputs[idx] <- data:
				counter++
			case <-ctx.Done():
				return ctx.Err()
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			dataReceived := false
			for _, in := range inputs {
				select {
				case data, ok := <-in:
					if !ok {
						continue
					}
					dataReceived = true
					if strings.Contains(data, "no multiplexer") {
						continue
					}
					select {
					case output <- data:
					case <-ctx.Done():
						return ctx.Err()
					}
				default:
				}
			}
			if !dataReceived {
				select {
				case <-time.After(10 * time.Millisecond):
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
	}
}
