package engine

import (
	"errors"

	"github.com/maxperrimond/kurin/example/domain"
)

type (
	CreateUserRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
)

func (engine *exampleEngine) GetUser(id string) (*domain.User, error) {
	user := engine.userRepository.Get(id)
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (engine *exampleEngine) CreateUser(r *CreateUserRequest) (*domain.User, error) {
	// Validate the request

	user := &domain.User{
		Email:    r.Email,
		Username: r.Username,
	}
	user.GenerateId()

	engine.userRepository.Create(user)

	return user, nil
}

func (engine *exampleEngine) DeleteUser(id string) error {
	user := engine.userRepository.Get(id)
	if user == nil {
		return errors.New("user not found")
	}

	engine.userRepository.Delete(user)

	return nil
}

func (engine *exampleEngine) ListUsers() []*domain.User {
	return engine.userRepository.List()
}
