package adapters

import (
	"net/http"

	"log"

	"fmt"
	"time"

	"github.com/maxperrimond/kurin"
)

type (
	HTTPAdapter struct {
		srv     *http.Server
		port    string
		Healthy bool
	}
)

func NewHTTPAdapter(handler http.Handler, port string) kurin.Adapter {
	a := &HTTPAdapter{
		port:    port,
		Healthy: true,
	}

	a.srv = &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return a
}

func (a *HTTPAdapter) Open() {
	log.Printf("Listening on http://0.0.0.0:%s\n", a.port)
	a.srv.ListenAndServe()
}

func (a *HTTPAdapter) Close() {
	err := a.srv.Close()
	if err != nil {
		log.Println(err)
	}
}

func (a *HTTPAdapter) ListenFailure(ce <-chan error) {
	go func() {
		err := <-ce
		if err != nil {
			a.Healthy = false
		}
	}()
}
