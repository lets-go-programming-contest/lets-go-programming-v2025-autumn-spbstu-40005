package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

const (
	prefixMarker        = "decorated: "
	noDecoratorMarker   = "no decorator"
	noMultiplexerMarker = "no multiplexer"
)

var (
	errDecorationRejected     = errors.New("can't be decorated")
	errNoSourcesProvided      = errors.New("no source channels provided")
	errNoDestinationsProvided = errors.New("no destination channels provided")
)

func PrefixDecoratorFunc(ctx context.Context, source chan string, destination chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case content, ok := <-source:
			if !ok {
				return nil
			}

			if strings.Contains(content, noDecoratorMarker) {
				return errDecorationRejected
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
		return errNoDestinationsProvided
	}

	distributionIndex := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case item, ok := <-source:
			if !ok {
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
	if len(sources) == 0 {
		return errNoSourcesProvided
	}

	var workerGroup sync.WaitGroup

	for _, src := range sources {
		workerGroup.Add(1)

		go func(src chan string) {
			defer workerGroup.Done()

/*			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-src:
					if !ok {
						return
					}

					if strings.Contains(data, noMultiplexerMarker) {
						continue
					}

					select {
					case destination <- data:
					case <-ctx.Done():
						return
					}
				}
			}*/
			for data := range src { // Выйдет сам при закрытии src
				if strings.Contains(data, noMultiplexerMarker) {
					continue
				}

				select {
				case destination <- data:
				case <-ctx.Done():
					return
				}
			}
		}(src)
	}

	workerGroup.Wait()

	return nil
}
