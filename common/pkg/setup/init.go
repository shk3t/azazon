package setup

import (
	"os"
	"os/signal"
	"sync"
)

type InitFunc func(workDir string) error
type DeinitFunc func()

type Initializer struct {
	init          InitFunc
	deinit        DeinitFunc
	up            bool
	mutex         sync.Mutex
	interruptChan chan os.Signal
}

func CreateInitFuncs(init InitFunc, deinit DeinitFunc) (InitFunc, DeinitFunc) {
	i := Initializer{
		init:          init,
		deinit:        deinit,
		up:            false,
		interruptChan: make(chan os.Signal, 1),
	}
	return i.Init, i.Deinit
}

func (i *Initializer) Init(workDir string) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if i.up {
		return nil
	}
	i.up = true

	if err := i.init(workDir); err != nil {
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