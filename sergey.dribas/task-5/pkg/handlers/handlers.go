package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

const (
	noFoundDecorator string = "no decorator"
	decorate         string = "decorated: "
	noFoundMultiplex string = "no multiplexer"
)

var (
	ErrNoDecorated   = errors.New("can't be decorated")
	ErrNoInputChanel = errors.New("no input channels provided")
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)

	for {
		select {
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, noFoundDecorator) {
				return ErrNoDecorated
			}

			if !strings.HasPrefix(data, decorate) {
				data = decorate + data
			}
			select {
			case output <- data:
			case <-ctx.Done():
				return nil
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, ch := range outputs {
			close(ch)
		}
	}()

	if len(outputs) == 0 {
		return nil
	}

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
				return nil
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	if len(inputs) == 0 {
		<-ctx.Done()

		return ErrNoInputChanel
	}

	var wait sync.WaitGroup

	wait.Add(len(inputs))

	for _, input := range inputs {
		reader := func() {
			defer wait.Done()

			for {
				select {
				case data, ok := <-input:
					if !ok {
						return
					}

					if strings.Contains(data, noFoundMultiplex) {
						continue
					}
					select {
					case output <- data:
					case <-ctx.Done():
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}

		go reader()
	}

	done := make(chan struct{})
	go func() {
		wait.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return nil
	}
}
