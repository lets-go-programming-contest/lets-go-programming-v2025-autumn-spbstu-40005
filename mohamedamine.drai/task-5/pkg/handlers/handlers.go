package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrDecorateBlocked = errors.New("can't be decorated")
	ErrNoTargets       = errors.New("no output channels provided")
)

const (
	skipDecorateKey = "no decorator"
	skipMuxKey      = "no multiplexer"
	prefix          = "decorated: "
)

func DecoratePrefix(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case s, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(s, skipDecorateKey) {
				return ErrDecorateBlocked
			}

			if !strings.HasPrefix(s, prefix) {
				s = prefix + s
			}

			select {
			case output <- s:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MergeInputs(ctx context.Context, ins []chan string, out chan string) error {
	var wg sync.WaitGroup
	wg.Add(len(ins))

	for _, c := range ins {
		ch := c

		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case v, ok := <-ch:
					if !ok {
						return
					}

					if strings.Contains(v, skipMuxKey) {
						continue
					}

					select {
					case out <- v:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	wg.Wait()
	return nil
}

func SplitStream(ctx context.Context, input chan string, outs []chan string) error {
	if len(outs) == 0 {
		return ErrNoTargets
	}

	pos := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case v, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outs[pos] <- v:
			case <-ctx.Done():
				return nil
			}

			pos = (pos + 1) % len(outs)
		}
	}
}
