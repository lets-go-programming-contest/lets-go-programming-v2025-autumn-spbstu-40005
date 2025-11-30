package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

const (
	tagDecorated      = "decorated: "
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
	var waitGrp sync.WaitGroup

	waitGrp.Add(len(inputs))

	for _, channel := range inputs {
		go func(inputChannel chan string) {
			defer waitGrp.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case message, isOpen := <-inputChannel:
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

	waitGrp.Wait()

	return nil
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	if len(outputs) == 0 {
		return nil
	}

	var assignmentIndex int

	for {
		select {
		case <-ctx.Done():
			return nil
		case message, isOpen := <-input:
			if !isOpen {
				return nil
			}

			targetChannel := outputs[assignmentIndex%len(outputs)]
			assignmentIndex++

			select {
			case targetChannel <- message:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
