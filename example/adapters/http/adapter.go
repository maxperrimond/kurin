package http

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/maxperrimond/kurin"
	httpAdapter "github.com/maxperrimond/kurin/adapters/http"
	"github.com/maxperrimond/kurin/example/engine"
	"go.uber.org/zap"
)

func NewHTTPAdapter(e engine.Engine, port int, logger *zap.Logger) kurin.Adapter {
	r := mux.NewRouter().StrictSlash(false)
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

	return httpAdapter.NewHTTPAdapter(h, port, "1.0.0", logger.Sugar())
}
