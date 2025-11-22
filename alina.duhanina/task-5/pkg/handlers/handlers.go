package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrDecorationRejected = errors.New("can't be decorated")
)

func PrefixDecoratorFunc(ctx context.Context, source chan string, destination chan string) error {
	defer close(destination)
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
	defer func() {
		for _, dest := range destinations {
			close(dest)
		}
	}()

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
	defer close(destination)

	if len(sources) == 0 {
		return nil
	}

	var wg sync.WaitGroup

	for _, src := range sources {
		wg.Add(1)
		go func(input chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, active := <-input:
					if !active {
						return
					}

					if strings.Contains(data, "no multiplexer") {
						continue
					}

					select {
					case destination <- data:
					case <-ctx.Done():
						return
					}
				}
			}
		}(src)
	}

	wg.Wait()
	return nil
}
