package http

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/maxperrimond/kurin"
	"github.com/maxperrimond/kurin/example/engine"
)

type (
	adapter struct {
		srv     *http.Server
		port    string
		healthy bool
	}
)

func NewHTTPAdapter(e engine.Engine, port string) kurin.Adapter {
	a := &adapter{
		port:    port,
		healthy: true,
	}

	r := mux.NewRouter().StrictSlash(false)

	r.NewRoute().
		Name("Health check").
		Methods(http.MethodGet).
		Path("/health").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if a.healthy {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusServiceUnavailable)
			}
		})

	r.NewRoute().
		Name("List all users").
		Methods(http.MethodGet).
		Path("/users").
		Handler(listUsersHandler(e))
	r.NewRoute().
		Name("Create user").
		Methods(http.MethodPost).
		Path("/users").
		Handler(createUserHandler(e))
	r.NewRoute().
		Name("Get user").
		Methods(http.MethodGet).
		Path("/users/{id}").
		Handler(getUserHandler(e))
	r.NewRoute().
		Name("Delete user").
		Methods(http.MethodDelete).
		Path("/users/{id}").
		Handler(deleteUserHandler(e))

	h := handlers.RecoveryHandler()(r)
	h = handlers.CompressHandler(h)
	h = handlers.ContentTypeHandler(h, "application/json")
	h = handlers.CombinedLoggingHandler(os.Stdout, h)

	a.srv = &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      h,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return a
}

func (a *adapter) Open() {
	log.Printf("Listening on http://0.0.0.0:%s\n", a.port)
	err := a.srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (a *adapter) Close() {
	err := a.srv.Close()
	if err != nil {
		log.Println(err)
	}
}

func (a *adapter) ListenFailure(err <-chan error) {
	go func() {
		newError := <-err
		if newError != nil {
			a.healthy = false
		}
	}()
}
