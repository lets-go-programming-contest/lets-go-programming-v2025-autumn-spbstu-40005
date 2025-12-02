package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

const (
	prefixToAdd   = "decorated: "
	triggerNoDeco = "no decorator"
	triggerNoMult = "no multiplexer"
)

var (
	ErrCantDecorate = errors.New("can't be decorated")
	ErrEmptyOutputs = errors.New("empty outputs")
)

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, triggerNoDeco) {
				return ErrCantDecorate
			}

			if !strings.HasPrefix(data, prefixToAdd) {
				data = prefixToAdd + data
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return ErrEmptyOutputs
	}

	defer close(output)

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	doneCh := make(chan struct{})

	for _, inCh := range inputs {
		go func(ch chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-ch:
					if !ok {
						return
					}

					if strings.Contains(data, triggerNoMult) {
						continue
					}

					select {
					case output <- data:
					case <-ctx.Done():
						return
					case <-doneCh:
						return
					}
				}
			}
		}(inCh)
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-ctx.Done():
		return nil
	case <-doneCh:
		return nil
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrEmptyOutputs
	}

	defer func() {
		for _, ch := range outputs {
			close(ch)
		}
	}()

	currentIndex := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			targetChan := outputs[currentIndex]
			currentIndex = (currentIndex + 1) % len(outputs)

			select {
			case targetChan <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
