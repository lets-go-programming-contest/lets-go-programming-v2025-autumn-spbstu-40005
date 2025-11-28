package handlers

import (
	"context"
	"errors"
	"strings"
)

const (
	decorator_prefix    = "decorated: "
	no_decorator_prefix = "no decorator"
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
