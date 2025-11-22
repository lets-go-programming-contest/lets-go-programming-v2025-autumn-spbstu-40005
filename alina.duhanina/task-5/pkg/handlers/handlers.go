package handlers

import (
	"context"
	"errors"
	"strings"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	prefix := "decorated: "

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return errors.New("can't be decorated")
			}

			if !strings.HasPrefix(data, prefix) {
				data = prefix + data
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	counter := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			idx := counter % len(outputs)
			select {
			case outputs[idx] <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
			counter++
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-inputs[0]:
			if !ok {
				inputs[0] = nil
			} else if !strings.Contains(data, "no multiplexer") {
				output <- data
			}
		case data, ok := <-inputs[1]:
			if !ok {
				inputs[1] = nil
			} else if !strings.Contains(data, "no multiplexer") {
				output <- data
			}
		}

		if inputs[0] == nil && inputs[1] == nil {
			return nil
		}
	}
}
