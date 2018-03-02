package kurin

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type (
	App struct {
		name            string
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
		Healthy() bool
		ListenFailure(<-chan error)
	}
)

func NewApp(name string, adapters ...Adapter) *App {
	closableSystems := make([]Closable, len(adapters))
	for i, a := range adapters {
		closableSystems[i] = Closable(a)
	}

	var fallibleSystems []Fallible
	for _, a := range adapters {
		f, ok := a.(Fallible)
		if ok {
			fallibleSystems = append(fallibleSystems, f)
		}
	}

	return &App{
		name:            name,
		adapters:        adapters,
		closableSystems: closableSystems,
		fallibleSystems: fallibleSystems,
	}
}

func (a *App) RegisterFallibleSystems(systems ...Fallible) {
	a.fallibleSystems = append(a.fallibleSystems, systems...)
}

func (a *App) RegisterClosableSystems(systems ...Closable) {
	a.closableSystems = append(a.closableSystems, systems...)
}

func (a *App) Run() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	defer close(stop)

	log.Printf("Starting %s application...\n", a.name)

	a.fail = make(chan error)
	defer close(a.fail)

	for _, system := range a.fallibleSystems {
		system.NotifyFail(a.fail)
	}

	for _, adapter := range a.adapters {
		go adapter.Open()
		adapter.ListenFailure(a.fail)
	}

	<-stop

	log.Println("Shutdown signal received, exiting...")

	for _, c := range a.closableSystems {
		c.Close()
	}
}
