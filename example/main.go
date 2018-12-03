package main

import (
	"github.com/maxperrimond/kurin"
	"github.com/maxperrimond/kurin/example/adapters/http"
	"github.com/maxperrimond/kurin/example/engine"
	"github.com/maxperrimond/kurin/example/providers/example"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()

	// Providers
	exampleProviderFactory := example.NewFactory()

	// Engine
	engineFactory := engine.NewFactory(exampleProviderFactory)
	e := engineFactory.NewEngine()

	// App
	a := kurin.NewApp("Example", http.NewHTTPAdapter(e, 7272, logger))
	a.RegisterSystems(exampleProviderFactory)
	a.Run()
}
