package http

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"context"

	"github.com/maxperrimond/kurin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type (
	Adapter struct {
		srv       *http.Server
		port      string
		version   string
		healthy   bool
		lastError error
	}
)

func NewHTTPAdapter(handler http.Handler, port string, version string) kurin.Adapter {
	adapter := &Adapter{
		port:    port,
		version: version,
		healthy: true,
	}

	totalCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_requests_total",
			Help: "A counter for requests to the wrapped handler.",
		},
		[]string{"code", "method"},
	)
	durationHist := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "app_response_duration_seconds",
			Help:        "A histogram of request latencies.",
			Buckets:     prometheus.DefBuckets,
			ConstLabels: prometheus.Labels{"handler": "api"},
		},
		[]string{"code", "method"},
	)
	prometheus.MustRegister(totalCount, durationHist)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if adapter.healthy {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(adapter.lastError.Error()))
		}
	})
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, version)
	})
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/",
		promhttp.InstrumentHandlerCounter(totalCount,
			promhttp.InstrumentHandlerDuration(durationHist, handler)))

	adapter.srv = &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return adapter
}

func (adapter *Adapter) Open() {
	log.Printf("Listening on http://0.0.0.0:%s\n", adapter.port)
	adapter.srv.ListenAndServe()
}

func (adapter *Adapter) Close() {
	err := adapter.srv.Shutdown(context.Background())
	if err != nil {
		log.Println(err)
	}
}

func (adapter *Adapter) OnFailure(err error) {
	if err != nil {
		adapter.lastError = err
		adapter.healthy = false
	}
}
