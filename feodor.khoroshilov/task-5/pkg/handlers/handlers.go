package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrCantBeDecorated = errors.New("can't be decorated")
)

const (
	skipDecorator   = "no decorator"
	decoration      = "decorated:"
	skipMultiplexer = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}
			
			if strings.Contains(data, skipDecorator) {
				return ErrCantBeDecorated
			}
			
			if !strings.HasPrefix(data, decoration) {
				data = decoration + data
			}
			
			select {
			case <-ctx.Done():
				return nil
			case output <- data:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}
	
	counter := 0
	
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}
			
			outputIndex := counter % len(outputs)
			counter++
			
			select {
			case <-ctx.Done():
				return nil
			case outputs[outputIndex] <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var waitGroup sync.WaitGroup

	for _, inputChan := range ins {
		waitGroup.Add(1)

		go func(inputCh chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-inputCh:
					if !ok {
						return
					}

					if strings.Contains(data, skipMultiplexer) {
						continue
					}

					out <- data
				}
			}
		}(inputChan)
	}

	waitGroup.Wait()

	return nil
}