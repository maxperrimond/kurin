package kurin

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type (
	App struct {
		name            string
		logger          Logger
		adapters        []Adapter
		fallibleSystems []Fallible
		closableSystems []Closable
		fail            chan error
	}

	Fallible interface {
		NotifyFail(chan error)
	}

	Closable interface {
		Close()
	}

	Adapter interface {
		Closable
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
	}
}

func (a *App) Run() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	defer close(stop)

	a.logger.Info(fmt.Sprintf("Starting %s application...", a.name))

	a.fail = make(chan error)
	defer close(a.fail)

	for _, system := range a.fallibleSystems {
		system.NotifyFail(a.fail)
	}

	for _, adapter := range a.adapters {
		go adapter.Open()
	}

	func() {
		for {
			select {
			case err := <-a.fail:
				for _, adapter := range a.adapters {
					adapter.OnFailure(err)
				}
				break
			case <-stop:
				return
			}
		}
	}()

	a.logger.Info("Shutdown signal received, exiting...")

	for _, c := range a.closableSystems {
		c.Close()
	}
}
