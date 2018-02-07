package engine

import "github.com/maxperrimond/kurin/example/domain"

type (
	Engine interface {
		GetUser(id string) (*domain.User, error)
		CreateUser(r *CreateUserRequest) (*domain.User, error)
		DeleteUser(id string) error
		ListUsers() []*domain.User
	}

	exampleEngine struct {
		userRepository UserRepository
	}
)
