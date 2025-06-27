package setup

import (
	"os"
	"os/signal"
	"sync"
)

type initFunc func(...any) error
type deinitFunc func()

type Initializer struct {
	init          initFunc
	deinit        deinitFunc
	up            bool
	mutex         sync.Mutex
	interruptChan chan os.Signal
}

func NewInitializer(init initFunc, deinit deinitFunc) *Initializer {
	return &Initializer{
		init:          init,
		deinit:        deinit,
		up:            false,
		interruptChan: make(chan os.Signal, 1),
	}
}

func (i *Initializer) Init(args ...any) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if i.up {
		return nil
	}
	i.up = true

	if err := i.init(args...); err != nil {
		i.up = false
		i.deinit()
		return err
	}

	signal.Notify(i.interruptChan, os.Interrupt)
	go func() {
		_, ok := <-i.interruptChan
		if ok {
			i.deinit()
			os.Exit(0)
		}
	}()
	return nil
}

func (i *Initializer) Deinit() {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if !i.up {
		return
	}
	i.up = false

	i.deinit()

	signal.Stop(i.interruptChan)
	close(i.interruptChan)
}