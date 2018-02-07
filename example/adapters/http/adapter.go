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
		srv  *http.Server
		port string
	}
)

func NewHTTPAdapter(e engine.Engine, port string) kurin.Adapter {
	r := mux.NewRouter()

	r.NewRoute().Name("List all users").Methods(http.MethodGet).Path("/users").Handler(listUsersHandler(e))
	r.NewRoute().Name("Create user").Methods(http.MethodPost).Path("/users").Handler(createUserHandler(e))
	r.NewRoute().Name("Get user").Methods(http.MethodGet).Path("/users/{id}").Handler(getUserHandler(e))
	r.NewRoute().Name("Delete user").Methods(http.MethodDelete).Path("/users/{id}").Handler(deleteUserHandler(e))

	h := handlers.RecoveryHandler()(r)
	h = handlers.CompressHandler(h)
	h = handlers.ContentTypeHandler(h, "application/json")
	h = handlers.CombinedLoggingHandler(os.Stdout, h)

	return &adapter{
		port: port,
		srv: &http.Server{
			Addr:         fmt.Sprintf(":%s", port),
			Handler:      h,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
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
