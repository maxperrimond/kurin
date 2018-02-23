package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/maxperrimond/kurin/example/engine"
)

func listUsersHandler(e engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := e.ListUsers()

		j, err := json.Marshal(users)
		if err != nil {
			panic(err)
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	}
}

func createUserHandler(e engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cur := &engine.CreateUserRequest{}
		json.NewDecoder(r.Body).Decode(cur)

		user, err := e.CreateUser(cur)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		j, err := json.Marshal(user)
		if err != nil {
			panic(err)
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(j)
	}
}

func getUserHandler(e engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		user, err := e.GetUser(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		j, err := json.Marshal(user)
		if err != nil {
			panic(err)
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)
	}
}

func deleteUserHandler(e engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		err := e.DeleteUser(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
