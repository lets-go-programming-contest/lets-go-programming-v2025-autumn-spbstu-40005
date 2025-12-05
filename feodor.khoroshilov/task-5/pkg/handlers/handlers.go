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
	var wg sync.WaitGroup
	var once sync.Once
	wg.Add(len(inputs))
	
	errCh := make(chan error, 1)
	
	defer func() {
		once.Do(func() {
		})
	}()
	
	for _, input := range inputs {
		go func(in chan string) {
			defer wg.Done()
			
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-in:
					if !ok {
						return
					}
					
					if strings.Contains(data, skipMultiplexer) {
						continue
					}
					
					select {
					case <-ctx.Done():
						return
					default:
						select {
						case <-ctx.Done():
							return
						case output <- data:
						default:
							select {
							case <-ctx.Done():
								return
							case output <- data:
							}
						}
					}
				}
			}
		}(input)
	}
	
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()
	
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}