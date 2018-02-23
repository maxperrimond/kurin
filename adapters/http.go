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
		healthy bool
	}
)

func NewHTTPAdapter(handler http.Handler, port string) kurin.Adapter {
	adapter := &HTTPAdapter{
		port:    port,
		healthy: true,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if adapter.healthy {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})
	mux.Handle("/", handler)

	adapter.srv = &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return adapter
}

func (adapter *HTTPAdapter) Open() {
	log.Printf("Listening on http://0.0.0.0:%s\n", adapter.port)
	adapter.srv.ListenAndServe()
}

func (adapter *HTTPAdapter) Close() {
	err := adapter.srv.Close()
	if err != nil {
		log.Println(err)
	}
}

func (adapter *HTTPAdapter) Healthy() bool {
	return adapter.healthy
}

func (adapter *HTTPAdapter) ListenFailure(ce <-chan error) {
	go func() {
		err := <-ce
		if err != nil {
			adapter.healthy = false
		}
	}()
}
