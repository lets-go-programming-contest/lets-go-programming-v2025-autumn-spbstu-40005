package handlers

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrDecorationRejected = errors.New("can't be decorated")
)

func PrefixDecoratorFunc(ctx context.Context, source chan string, destination chan string) error {
	const prefixMarker = "decorated: "

	for {
		select {
		case <-ctx.Done():
			return nil
		case content, active := <-source:
			if !active {
				return nil
			}

			if strings.Contains(content, "no decorator") {
				return ErrDecorationRejected
			}

			if !strings.HasPrefix(content, prefixMarker) {
				content = prefixMarker + content
			}

			select {
			case destination <- content:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, source chan string, destinations []chan string) error {
	if len(destinations) == 0 {
		return nil
	}

	distributionIndex := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case item, active := <-source:
			if !active {
				return nil
			}

			target := destinations[distributionIndex%len(destinations)]
			distributionIndex++

			select {
			case target <- item:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, sources []chan string, destination chan string) error {
	type workerResult struct {
		data string
		done bool
	}

	results := make(chan workerResult, len(sources))

	for _, src := range sources {
		go func(input chan string) {
			for {
				select {
				case <-ctx.Done():
					return
				case data, active := <-input:
					if !active {
						results <- workerResult{done: true}
						return
					}

					if strings.Contains(data, "no multiplexer") {
						continue
					}

					select {
					case results <- workerResult{data: data}:
					case <-ctx.Done():
						return
					}
				}
			}
		}(src)
	}

	activeWorkers := len(sources)
	for activeWorkers > 0 {
		select {
		case <-ctx.Done():
			return nil
		case result := <-results:
			if result.done {
				activeWorkers--
			} else {
				select {
				case destination <- result.data:
				case <-ctx.Done():
					return nil
				}
			}
		}
	}

	return nil
}
