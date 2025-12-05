package handlers

import (
	"context"
	"strings"
	"sync"
)

const (
	PrefixDecoratorMarker = "decorated: "
	NoDecoratorToken      = "no decorator"
	NoMultiplexerToken    = "no multiplexer"
)

var (
	ErrRejectedByDecorator = NewHandlerError("rejected by decorator")
	ErrEmptyOutputs        = NewHandlerError("outputs must not be empty")
)

type HandlerError struct {
	msg string
}

func NewHandlerError(msg string) error { return &HandlerError{msg: msg} }
func (e *HandlerError) Error() string  { return e.msg }

func PrefixDecorator(ctx context.Context, in <-chan string, out chan<- string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-in:
			if !ok {
				return nil
			}
			if strings.Contains(data, NoDecoratorToken) {
				return ErrRejectedByDecorator
			}
			if !strings.HasPrefix(data, PrefixDecoratorMarker) {
				data = PrefixDecoratorMarker + data
			}
			select {
			case out <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func RoundRobinSeparator(ctx context.Context, in <-chan string, outs []chan<- string) error {
	if len(outs) == 0 {
		return ErrEmptyOutputs
	}

	i := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-in:
			if !ok {
				return nil
			}
			ch := outs[i%len(outs)]
			i++
			select {
			case ch <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func Multiplexer(ctx context.Context, ins []<-chan string, out chan<- string) error {
	type source struct {
		ch <-chan string
	}

	sources := make([]source, len(ins))
	for i, ch := range ins {
		sources[i] = source{ch: ch}
	}

	var wg sync.WaitGroup
	for _, src := range sources {
		wg.Add(1)
		go func(ch <-chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-ch:
					if !ok {
						return
					}
					if strings.Contains(data, NoMultiplexerToken) {
						continue
					}
					select {
					case out <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(src.ch)
	}

	wg.Wait()
	return nil
}
