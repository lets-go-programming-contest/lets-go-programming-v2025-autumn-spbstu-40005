package handlers

import (
	"context"
	"errors"
	"strings"
)

const (
	tagDecorated      = "decorated"
	stopDecorator     = "no decorator"
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
