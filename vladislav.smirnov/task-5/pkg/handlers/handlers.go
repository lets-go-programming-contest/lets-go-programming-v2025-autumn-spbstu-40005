package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	ErrCantDecorate = errors.New("can't be decorated")
	ErrEmptyChannel = errors.New("empty channel")
)

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	const prefix = "decorated: "

	for {
		select {
		case <-ctx.Done():
			return nil

		case value, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(value, "no decorator") {
				return ErrCantDecorate
			}

			if !strings.HasPrefix(value, prefix) {
				value = prefix + value
			}

			select {
			case output <- value:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrEmptyChannel
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case value, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outputs[index] <- value:
				index = (index + 1) % len(outputs)

			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return ErrEmptyChannel
	}

	var waitGroup sync.WaitGroup

	errGroup, gCtx := errgroup.WithContext(ctx)

	for _, channel := range inputs {
		waitGroup.Add(1)
		errGroup.Go(func(c chan string) func() error {
			defer waitGroup.Done()
			return func() error {
				for {
					select {
					case <-gCtx.Done():
						return gCtx.Err()

					case value, ok := <-channel:
						if !ok {
							return nil
						}

						if strings.Contains(value, "no multiplexer") {
							continue
						}

						select {
						case output <- value:

						case <-gCtx.Done():
							return gCtx.Err()
						}
					}
				}
			}
		}(channel))
	}

	err := errGroup.Wait()

	if err != nil && errors.Is(err, context.Canceled) {
		return nil
	}

	return fmt.Errorf("multiplexer failed: %w", err)
}
