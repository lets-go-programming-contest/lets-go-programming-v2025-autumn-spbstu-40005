package handlers

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync/atomic"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	const prefix = "decorated: "
	const errorSubstring = "no decorator"
	const errorMsg = "can’t be decorated"

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, errorSubstring) {
				return errors.New(errorMsg)
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
	if len(outputs) == 0 {
		return nil
	}

	var counter int64 = -1

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			idx := atomic.AddInt64(&counter, 1) % int64(len(outputs))
			select {
			case outputs[idx] <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	const skipSubstring = "no multiplexer"

	if len(inputs) == 0 {
		return nil
	}

	cases := make([]reflect.SelectCase, len(inputs)+1)
	for i, ch := range inputs {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
	}
	cases[len(inputs)] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ctx.Done()),
	}

	for {
		chosen, value, ok := reflect.Select(cases)
		if chosen == len(inputs) {
			return ctx.Err()
		}

		if !ok {
			return nil
		}

		data := value.String()
		if strings.Contains(data, skipSubstring) {
			continue
		}

		select {
		case output <- data:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
