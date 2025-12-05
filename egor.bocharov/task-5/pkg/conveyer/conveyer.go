package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const undefined = "undefined"

var ErrChannelNotFound = errors.New("chan not found")

type conveyerImpl struct {
	size     int
	channels map[string]chan string
	handlers []func(ctx context.Context) error
	mu       sync.RWMutex
}

// New создает новый конвейер с указанным размером буфера каналов.
func New(size int) *conveyerImpl {
	return &conveyerImpl{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]func(ctx context.Context) error, 0),
		mu:       sync.RWMutex{},
	}
}

// RegisterDecorator регистрирует обработчик-модификатор данных.
func (c *conveyerImpl) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input, output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputCh := c.getOrCreateChannel(input)
	outputCh := c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputCh, outputCh)
	})
}

// RegisterMultiplexer регистрирует мультиплексор.
func (c *conveyerImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	outputCh := c.getOrCreateChannel(output)
	inputChs := make([]chan string, len(inputs))

	for i, name := range inputs {
		inputChs[i] = c.getOrCreateChannel(name)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputChs, outputCh)
	})
}

// RegisterSeparator регистрирует сепаратор.
func (c *conveyerImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputCh := c.getOrCreateChannel(input)
	outputChs := make([]chan string, len(outputs))

	for i, name := range outputs {
		outputChs[i] = c.getOrCreateChannel(name)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputCh, outputChs)
	})
}

// Run запускает все зарегистрированные обработчики в отдельных горутинах.
// Блокируется до завершения всех обработчиков или отмены контекста.
func (c *conveyerImpl) Run(ctx context.Context) error {
	defer c.closeAllChannels()

	errgr, ctx := errgroup.WithContext(ctx)

	for _, h := range c.handlers {
		handler := h // capture loop variable
		errgr.Go(func() error {
			return handler(ctx)
		})
	}

	if err := errgr.Wait(); err != nil {
		return fmt.Errorf("conveyer run failed: %w", err)
	}

	return nil
}

// closeAllChannels закрывает все каналы конвейера.
func (c *conveyerImpl) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, ch := range c.channels {
		close(ch)
	}
}

// Send отправляет данные в канал с указанным идентификатором.
func (c *conveyerImpl) Send(ctx context.Context, input string, data string) error {
	c.mu.RLock()
	ch, exists := c.channels[input]
	c.mu.RUnlock()

	if !exists {
		return ErrChannelNotFound
	}

	select {
	case ch <- data:
		return nil
	case <-ctx.Done(): // на всякий случай, если контекст уже отменен
		return ctx.Err()
	}
}

// Recv получает данные из канала с указанным идентификатором.
func (c *conveyerImpl) Recv(ctx context.Context, output string) (string, error) {
	c.mu.RLock()
	ch, exists := c.channels[output]
	c.mu.RUnlock()

	if !exists {
		return "", ErrChannelNotFound
	}

	select {
	case val, ok := <-ch:
		if !ok {
			return undefined, nil
		}
		return val, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// getOrCreateChannel возвращает существующий канал или создает новый.
func (c *conveyerImpl) getOrCreateChannel(name string) chan string {
	c.mu.RLock()
	if ch, exists := c.channels[name]; exists {
		c.mu.RUnlock()
		return ch
	}
	c.mu.RUnlock()

	// Создаем новый канал под блокировкой
	c.mu.Lock()
	defer c.mu.Unlock()
	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}
