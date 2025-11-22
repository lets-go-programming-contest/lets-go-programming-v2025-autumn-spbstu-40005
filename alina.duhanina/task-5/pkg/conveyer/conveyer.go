package conveyer

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound = errors.New("chan not found")
	UndefinedValue  = "undefined"
)

type DataProcessor interface {
	Process(ctx context.Context) error
}

type decoratorTask struct {
	fn     func(ctx context.Context, in, out chan string) error
	input  chan string
	output chan string
}

func (d *decoratorTask) Process(ctx context.Context) error {
	return d.fn(ctx, d.input, d.output)
}

type multiplexerTask struct {
	fn      func(ctx context.Context, ins []chan string, out chan string) error
	inputs  []chan string
	output  chan string
}

func (m *multiplexerTask) Process(ctx context.Context) error {
	return m.fn(ctx, m.inputs, m.output)
}

type separatorTask struct {
	fn      func(ctx context.Context, in chan string, outs []chan string) error
	input   chan string
	outputs []chan string
}

func (s *separatorTask) Process(ctx context.Context) error {
	return s.fn(ctx, s.input, s.outputs)
}

type Conveyer struct {
	capacity int
	storage  *channelStorage
	tasks    []DataProcessor
	mu       sync.RWMutex
}

type channelStorage struct {
	channels map[string]chan string
	mu       sync.RWMutex
}

func newChannelStorage() *channelStorage {
	return &channelStorage{
		channels: make(map[string]chan string),
	}
}

func (cs *channelStorage) getOrCreate(name string, capacity int) chan string {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if ch, exists := cs.channels[name]; exists {
		return ch
	}

	ch := make(chan string, capacity)
	cs.channels[name] = ch
	return ch
}

func (cs *channelStorage) get(name string) (chan string, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	if ch, exists := cs.channels[name]; exists {
		return ch, nil
	}

	return nil, ErrChanNotFound
}

func (cs *channelStorage) closeAll() {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	for name, ch := range cs.channels {
		close(ch)
		delete(cs.channels, name)
	}
}

func New(capacity int) *Conveyer {
	return &Conveyer{
		capacity: capacity,
		storage:  newChannelStorage(),
		tasks:    make([]DataProcessor, 0),
	}
}

func (c *Conveyer) RegisterDecorator(
	fn func(ctx context.Context, input, output chan string) error,
	inputName, outputName string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputChan := c.storage.getOrCreate(inputName, c.capacity)
	outputChan := c.storage.getOrCreate(outputName, c.capacity)

	task := &decoratorTask{
		fn:     fn,
		input:  inputChan,
		output: outputChan,
	}

	c.tasks = append(c.tasks, task)
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputNames []string, outputName string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	outputChan := c.storage.getOrCreate(outputName, c.capacity)
	inputChans := make([]chan string, len(inputNames))

	for i, name := range inputNames {
		inputChans[i] = c.storage.getOrCreate(name, c.capacity)
	}

	task := &multiplexerTask{
		fn:     fn,
		inputs: inputChans,
		output: outputChan,
	}

	c.tasks = append(c.tasks, task)
}

func (c *Conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	inputName string, outputNames []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputChan := c.storage.getOrCreate(inputName, c.capacity)
	outputChans := make([]chan string, len(outputNames))

	for i, name := range outputNames {
		outputChans[i] = c.storage.getOrCreate(name, c.capacity)
	}

	task := &separatorTask{
		fn:      fn,
		input:   inputChan,
		outputs: outputChans,
	}

	c.tasks = append(c.tasks, task)
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.mu.Lock()
	defer func() {
		c.storage.closeAll()
		c.mu.Unlock()
	}()

	workerGroup, ctx := errgroup.WithContext(ctx)

	for _, task := range c.tasks {
		currentTask := task
		workerGroup.Go(func() error {
			return currentTask.Process(ctx)
		})
	}

	return workerGroup.Wait()
}

func (c *Conveyer) Send(channelName string, data string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, err := c.storage.get(channelName)
	if err != nil {
		return err
	}

	ch <- data
	return nil
}

func (c *Conveyer) Recv(channelName string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, err := c.storage.get(channelName)
	if err != nil {
		return "", err
	}

	value, active := <-ch
	if !active {
		return UndefinedValue, nil
	}

	return value, nil
}

func (c *Conveyer) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.storage.closeAll()
}
