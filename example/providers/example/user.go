package example

import (
	"fmt"
	"time"

	"github.com/maxperrimond/kurin/example/domain"
	"github.com/maxperrimond/kurin/example/engine"
)

type (
	userRepository struct {
		db userDB
	}
)

func newUserRepository(db userDB) engine.UserRepository {
	return &userRepository{db}
}

func (repository *userRepository) Get(id string) *domain.User {
	user, ok := repository.db[id]
	if !ok {
		return nil
	}

	return user
}

func (repository *userRepository) Create(user *domain.User) {
	_, ok := repository.db[user.Id]
	if ok {
		panic(fmt.Sprintf("User with id '%s' already exists", user.Id))
	}

	user.CreatedAt = time.Now()
	repository.db[user.Id] = user
}

func (repository *userRepository) Delete(user *domain.User) {
	delete(repository.db, user.Id)
}

func (repository *userRepository) List() []*domain.User {
	list := make([]*domain.User, len(repository.db))

	i := 0
	for _, user := range repository.db {
		list[i] = user
		i++
	}

	return list
}
