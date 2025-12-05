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
	wg       sync.WaitGroup // Добавляем WaitGroup для ожидания завершения
}

// New создаёт новый конвейер с указанным размером буфера каналов.
func New(size int) *conveyerImpl {
	return &conveyerImpl{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]func(ctx context.Context) error, 0),
		mu:       sync.RWMutex{},
		wg:       sync.WaitGroup{},
	}
}

// RegisterDecorator регистрирует обработчик-модификатор данных.
func (c *conveyerImpl) RegisterDecorator(
	decoratorFunc func(ctx context.Context, input chan string, output chan string) error,
	input, output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputChannel := c.getOrCreateChannel(input)
	outputChannel := c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return decoratorFunc(ctx, inputChannel, outputChannel)
	})
}

// RegisterMultiplexer регистрирует мультиплексор.
func (c *conveyerImpl) RegisterMultiplexer(
	multiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	outputChannel := c.getOrCreateChannel(output)
	inputChannels := make([]chan string, len(inputs))

	for i, name := range inputs {
		inputChannels[i] = c.getOrCreateChannel(name)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return multiplexerFunc(ctx, inputChannels, outputChannel)
	})
}

// RegisterSeparator регистрирует сепаратор.
func (c *conveyerImpl) RegisterSeparator(
	separatorFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputChannel := c.getOrCreateChannel(input)
	outputChannels := make([]chan string, len(outputs))

	for i, name := range outputs {
		outputChannels[i] = c.getOrCreateChannel(name)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return separatorFunc(ctx, inputChannel, outputChannels)
	})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	// Закрываем каналы после завершения всех обработчиков
	defer c.closeAllChannels()

	errGroup, ctx := errgroup.WithContext(ctx)

	for _, handler := range c.handlers {
		currentHandler := handler

		errGroup.Go(func() error {
			return currentHandler(ctx)
		})
	}

	// Ожидаем завершения всех обработчиков
	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("conveyer run failed: %w", err)
	}

	return nil
}

// closeAllChannels закрывает все каналы конвейера.
func (c *conveyerImpl) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for name, channel := range c.channels {
		// Проверяем, не закрыт ли уже канал
		select {
		case _, ok := <-channel:
			if !ok {
				continue // Канал уже закрыт
			}
		default:
		}

		// Закрываем канал
		close(channel)
		delete(c.channels, name)
	}
}

func (c *conveyerImpl) Send(input string, data string) error {
	c.mu.RLock()
	channel, exists := c.channels[input]
	c.mu.RUnlock()

	if !exists {
		return fmt.Errorf("conveyer send failed: %w", ErrChannelNotFound)
	}

	select {
	case channel <- data:
		return nil
	default:
		return fmt.Errorf("channel %s is full or closed", input)
	}
}

// Recv получает данные из канала с указанным идентификатором.
// Блокируется, пока данные не поступят или канал не закроется.
func (c *conveyerImpl) Recv(output string) (string, error) {
	c.mu.RLock()
	channel, exists := c.channels[output]
	c.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("conveyer recv failed: %w", ErrChannelNotFound)
	}

	val, ok := <-channel
	if !ok {
		return undefined, nil
	}

	return val, nil
}

func (c *conveyerImpl) getOrCreateChannel(name string) chan string {
	c.mu.RLock()
	if channel, exists := c.channels[name]; exists {
		c.mu.RUnlock()
		return channel
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Двойная проверка
	if channel, exists := c.channels[name]; exists {
		return channel
	}

	newChannel := make(chan string, c.size)
	c.channels[name] = newChannel

	return newChannel
}
