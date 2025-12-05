package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const Undefined = "undefined"

var (
	ErrChannelNotFound = errors.New("channel not found")
	ErrEmptyOutputs    = errors.New("outputs must not be empty")
)

// ChannelStore управляет каналами по имени.
type ChannelStore struct {
	size int
	mu   sync.RWMutex
	chs  map[string]chan string
}

func NewChannelStore(bufferSize int) *ChannelStore {
	return &ChannelStore{
		size: bufferSize,
		chs:  make(map[string]chan string),
	}
}

func (cs *ChannelStore) Get(name string) (chan string, bool) {
	cs.mu.RLock()
	ch, ok := cs.chs[name]
	cs.mu.RUnlock()
	return ch, ok
}

func (cs *ChannelStore) MustGet(name string) (chan string, error) {
	if ch, ok := cs.Get(name); ok {
		return ch, nil
	}
	return nil, ErrChannelNotFound
}

// GetOrCreate создаёт или возвращает существующий канал.
// Потокобезопасен.
func (cs *ChannelStore) GetOrCreate(name string) chan string {
	if ch, ok := cs.Get(name); ok {
		return ch
	}

	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Проверка второй раз под блокировкой (double-check idiom)
	if ch, ok := cs.chs[name]; ok {
		return ch
	}

	ch := make(chan string, cs.size)
	cs.chs[name] = ch
	return ch
}

func (cs *ChannelStore) CloseAll() {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	for _, ch := range cs.chs {
		close(ch)
	}
}

// Handler — любая функция, реализующая логику обработки данных в конвейере.
type Handler func(ctx context.Context) error

// Conveyer управляет жизненным циклом конвейера.
type Conveyer struct {
	store    *ChannelStore
	handlers []Handler
	mu       sync.Mutex // защищает handlers от 동시ной модификации
}

func New(bufferSize int) *Conveyer {
	return &Conveyer{
		store:    NewChannelStore(bufferSize),
		handlers: make([]Handler, 0),
	}
}

func (c *Conveyer) RegisterDecorator(fn func(context.Context, chan string, chan string) error, input, output string) {
	inCh := c.store.GetOrCreate(input)
	outCh := c.store.GetOrCreate(output)

	handler := func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	}

	c.mu.Lock()
	c.handlers = append(c.handlers, handler)
	c.mu.Unlock()
}

func (c *Conveyer) RegisterMultiplexer(fn func(context.Context, []chan string, chan string) error, inputs []string, output string) {
	outCh := c.store.GetOrCreate(output)
	inChs := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChs[i] = c.store.GetOrCreate(name)
	}

	handler := func(ctx context.Context) error {
		return fn(ctx, inChs, outCh)
	}

	c.mu.Lock()
	c.handlers = append(c.handlers, handler)
	c.mu.Unlock()
}

func (c *Conveyer) RegisterSeparator(fn func(context.Context, chan string, []chan string) error, input string, outputs []string) {
	inCh := c.store.GetOrCreate(input)
	outChs := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChs[i] = c.store.GetOrCreate(name)
	}

	handler := func(ctx context.Context) error {
		return fn(ctx, inCh, outChs)
	}

	c.mu.Lock()
	c.handlers = append(c.handlers, handler)
	c.mu.Unlock()
}

func (c *Conveyer) Run(ctx context.Context) error {
	defer c.store.CloseAll()

	// Копируем handlers под блокировкой — безопасно для конкурентного доступа
	var handlers []Handler
	c.mu.Lock()
	handlers = append(handlers, c.handlers...)
	c.mu.Unlock()

	eg, ctx := errgroup.WithContext(ctx)
	for _, h := range handlers {
		h := h // захват
		eg.Go(func() error { return h(ctx) })
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("conveyer failed: %w", err)
	}
	return nil
}

func (c *Conveyer) Send(channelName string, data string) error {
	ch, err := c.store.MustGet(channelName)
	if err != nil {
		return err
	}
	select {
	case ch <- data:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Conveyer) Recv(channelName string) (string, error) {
	ch, err := c.store.MustGet(channelName)
	if err != nil {
		return "", err
	}
	select {
	case val, ok := <-ch:
		if !ok {
			return Undefined, nil
		}
		return val, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
