package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

const (
	tagDecorated      = "decorated"
	stopDecorator     = "no decorator"
	stopMultiplexer   = "no multiplexer"
	msgCannotDecorate = "can't be decorated"
)

var ErrDecorationRefused = errors.New(msgCannotDecorate)

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case message, isOpen := <-input:
			if !isOpen {
				return nil
			}

			if strings.Contains(message, stopDecorator) {
				return ErrDecorationRefused
			}

			if !strings.HasPrefix(message, tagDecorated) {
				message = tagDecorated + message
			}

			select {
			case output <- message:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, channel := range inputs {
		go func(ch chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case message, isOpen := <-ch:
					if !isOpen {
						return
					}

					if strings.Contains(message, stopMultiplexer) {
						continue
					}

					select {
					case output <- message:
					case <-ctx.Done():
						return
					}
				}
			}
		}(channel)
	}

	wg.Wait()

	return nil
}
