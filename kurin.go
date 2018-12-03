package kurin

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type (
	App struct {
		name             string
		logger           Logger
		adapters         []Adapter
		fallibleSystems  []Fallible
		closableSystems  []Closable
		stoppableSystems []Stoppable
		stop             chan os.Signal
		fail             chan error
	}

	Fallible interface {
		NotifyFail(chan error)
	}

	Stoppable interface {
		NotifyStop(chan os.Signal)
	}

	Closable interface {
		Close()
	}

	Adapter interface {
		Closable
		Stoppable
		Open()
		OnFailure(error)
	}
)

func NewApp(name string, adapters ...Adapter) *App {
	app := &App{
		name:            name,
		logger:          &defaultLogger{},
		adapters:        adapters,
		closableSystems: make([]Closable, 0),
		fallibleSystems: make([]Fallible, 0),
	}
	app.RegisterSystems(adapters)

	return app
}

func (a *App) SetLogger(logger Logger) {
	a.logger = logger
}

func (a *App) RegisterSystems(systems ...interface{}) {
	for _, s := range systems {
		if f, ok := s.(Fallible); ok {
			a.fallibleSystems = append(a.fallibleSystems, f)
		}

		if c, ok := s.(Closable); ok {
			a.closableSystems = append(a.closableSystems, c)
		}

		if c, ok := s.(Stoppable); ok {
			c.NotifyStop(a.stop)
		}
	}
}

func (a *App) Run() {
	a.logger.Info(fmt.Sprintf("Starting %s application...", a.name))

	a.stop = make(chan os.Signal, 1)
	a.fail = make(chan error)
	defer close(a.fail)
	defer close(a.stop)

	signal.Notify(a.stop, syscall.SIGINT, syscall.SIGTERM)

	for _, system := range a.fallibleSystems {
		system.NotifyFail(a.fail)
	}

	for _, system := range a.stoppableSystems {
		system.NotifyStop(a.stop)
	}

	for _, adapter := range a.adapters {
		go adapter.Open()
	}

	func() {
		for {
			select {
			case err := <-a.fail:
				a.logger.Error(err)
				for _, adapter := range a.adapters {
					adapter.OnFailure(err)
				}
				break
			case <-a.stop:
				return
			}
		}
	}()

	a.logger.Info("Shutdown signal received, exiting...")

	for _, c := range a.closableSystems {
		c.Close()
	}
}
