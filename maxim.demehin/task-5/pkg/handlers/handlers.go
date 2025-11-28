package handlers

import (
	"context"
	"errors"
	"strings"
)

const (
	decorator_prefix    = "decorated: "
	no_decorator_prefix = "no decorator"
	no_multip_str = "no multiplexer"
)

var ErrCantDecorate = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, no_decorator_prefix) {
				return ErrCantDecorate
			}

			if !strings.HasPrefix(data, decorator_prefix) {
				data = decorator_prefix + data
			}

			output <- data

			return nil
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	out_len := len(outputs)

	if out_len == 0 {
		return nil
	}

	cnt := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			out_idx := cnt % out_len

			cnt++

			select {
			case <-ctx.Done():
				return nil
			case outputs[out_idx] <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case 
		}
	}
}