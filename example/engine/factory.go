package engine

type (
	Factory interface {
		NewEngine() Engine
	}

	engineFactory struct {
		ExampleProviderFactory
	}
)

func NewFactory(e ExampleProviderFactory) Factory {
	return &engineFactory{e}
}

func (f *engineFactory) NewEngine() Engine {
	return &exampleEngine{
		userRepository: f.NewUserRepository(),
	}
}
