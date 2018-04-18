package example

import (
	"github.com/maxperrimond/kurin/example/domain"
	"github.com/maxperrimond/kurin/example/engine"
)

type (
	userDB map[string]*domain.User

	ProviderFactory struct {
		db userDB
	}
)

func NewFactory() *ProviderFactory {
	return &ProviderFactory{
		db: userDB{},
	}
}

func (f *ProviderFactory) NewUserRepository() engine.UserRepository {
	return newUserRepository(f.db)
}

func (f *ProviderFactory) Close()                    {}
func (f *ProviderFactory) NotifyFail(err chan error) {}
