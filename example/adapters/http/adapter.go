package http

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/maxperrimond/kurin"
	"github.com/maxperrimond/kurin/example/engine"
)

func NewHTTPAdapter(e engine.Engine, port string) kurin.Adapter {
	r := mux.NewRouter().StrictSlash(false)
	h := handlers.RecoveryHandler()(r)
	h = handlers.CompressHandler(h)
	h = handlers.ContentTypeHandler(h, "application/json")
	h = handlers.CombinedLoggingHandler(os.Stdout, h)

	a := httpAdapter.NewHTTPAdapter(h, port)

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

	return a
}
