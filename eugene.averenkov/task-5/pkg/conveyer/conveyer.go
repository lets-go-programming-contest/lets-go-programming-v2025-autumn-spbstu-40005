package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

const undefined = "undefined"

type Conveyer struct {
	channels map[string]chan string
	size     int
	tasks    []func(ctx context.Context) error
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		size:     size,
		tasks:    []func(ctx context.Context) error{},
	}
}

func (c *Conveyer) getOrCreateChannel(name string) chan string {
	if ch, exists := c.channels[name]; exists {
		return ch
	}
	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inCh := c.getOrCreateChannel(input)
	outCh := c.getOrCreateChannel(output)

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return fn(ctx, inCh, outCh)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inChans := make([]chan string, len(inputs))
	for i, name := range inputs {
		inChans[i] = c.getOrCreateChannel(name)
	}
	outCh := c.getOrCreateChannel(output)

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return fn(ctx, inChans, outCh)
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inCh := c.getOrCreateChannel(input)
	outChans := make([]chan string, len(outputs))
	for i, name := range outputs {
		outChans[i] = c.getOrCreateChannel(name)
	}

	c.tasks = append(c.tasks, func(ctx context.Context) error {
		return fn(ctx, inCh, outChans)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
    var wg sync.WaitGroup
    errCh := make(chan error, len(c.tasks))

    for _, task := range c.tasks {
        wg.Add(1)
        go func(t func(ctx context.Context) error) {
            defer wg.Done()
            if err := t(ctx); err != nil {
                select {
                case errCh <- err:
                default:
                }
            }
        }(task)
    }

    go func() {
        wg.Wait()
        close(errCh)
    }()

    select {
    case err := <-errCh:
        return err
    case <-ctx.Done():
        return nil
    }
}

func (c *Conveyer) Send(input string, data string) error {
	ch, exists := c.channels[input]
	if !exists {
		return ErrChanNotFound
	}
	ch <- data
	return nil
}
func (c *Conveyer) Recv(output string) (string, error) {
    ch, exists := c.channels[output]
    if !exists {
        return "", ErrChanNotFound
    }

    data, ok := <-ch
    if !ok {
        return undefined, nil
    }
    return data, nil
}

/*package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	errChanNotFound    = errors.New("chan not found")
	errConveyerRunning = errors.New("conveyer is already running")
)

const errUndefined = "undefined"

type Conveyer struct {
	mu       sync.RWMutex
	channels map[string]chan string
	size     int
	tasks    []task
	running  bool
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

type task struct {
	fn      func(ctx context.Context) error
	inputs  []string
	outputs []string
}

func New(size int) *Conveyer {
	return &Conveyer{
		mu:       sync.RWMutex{},
		channels: make(map[string]chan string),
		size:     size,
		tasks:    []task{},
		running:  false,
		cancel:   nil,
		wg:       sync.WaitGroup{},
	}
}

func (c *Conveyer) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if channel, exists := c.channels[name]; exists {
		return channel
	}

	channel := make(chan string, c.size)
	c.channels[name] = channel

	return channel
}

func (c *Conveyer) RegisterDecorator(
	decoratorFunc func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inputChan := c.getOrCreateChannel(input)
	outputChan := c.getOrCreateChannel(output)

	taskFunc := func(ctx context.Context) error {
		return decoratorFunc(ctx, inputChan, outputChan)
	}

	c.tasks = append(c.tasks, task{
		fn:      taskFunc,
		inputs:  []string{input},
		outputs: []string{output},
	})
}

func (c *Conveyer) RegisterMultiplexer(
	multiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputChannels := make([]chan string, len(inputs))
	for i, inputName := range inputs {
		inputChannels[i] = c.getOrCreateChannel(inputName)
	}

	outputChannel := c.getOrCreateChannel(output)

	taskFunc := func(ctx context.Context) error {
		return multiplexerFunc(ctx, inputChannels, outputChannel)
	}

	c.tasks = append(c.tasks, task{
		fn:      taskFunc,
		inputs:  inputs,
		outputs: []string{output},
	})
}

func (c *Conveyer) RegisterSeparator(
	separatorFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inputChannel := c.getOrCreateChannel(input)

	outputChannels := make([]chan string, len(outputs))
	for i, outputName := range outputs {
		outputChannels[i] = c.getOrCreateChannel(outputName)
	}

	taskFunc := func(ctx context.Context) error {
		return separatorFunc(ctx, inputChannel, outputChannels)
	}

	c.tasks = append(c.tasks, task{
		fn:      taskFunc,
		inputs:  []string{input},
		outputs: outputs,
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()

		return errConveyerRunning
	}

	c.running = true
	runCtx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.running = false
		c.cancel = nil
		c.mu.Unlock()
	}()

	errorChannel := make(chan error, len(c.tasks))

	for _, currentTask := range c.tasks {
		c.wg.Add(1)

		taskCopy := currentTask
		go func() {
			defer c.wg.Done()

			if err := taskCopy.fn(runCtx); err != nil {
				select {
				case errorChannel <- err:
				default:
				}
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
		close(errorChannel)
	}()

	select {
	case <-done:
		return nil
	case err := <-errorChannel:
		cancel()

		return err
	case <-runCtx.Done():
		return nil
	}
}

func (c *Conveyer) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for name, channel := range c.channels {
		select {
		case <-channel:

		default:
		}
		close(channel)
		delete(c.channels, name)
	}
}

func (c *Conveyer) Send(input string, data string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	channel, exists := c.channels[input]
	if !exists {
		return errChanNotFound
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	channel, exists := c.channels[output]
	c.mu.RUnlock()

	if !exists {
		return "", errChanNotFound
	}

	data, ok := <- channel
	if !ok {
		return errUndefined, nil
	}

	return data, nil
}

func (c *Conveyer) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cancel != nil {
		c.cancel()
	}

	c.wg.Wait()
	c.closeAllChannels()
	c.running = false
}
*/
