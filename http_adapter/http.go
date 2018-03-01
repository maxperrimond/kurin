package http_adapter

import (
	"net/http"

	"log"

	"fmt"
	"time"

	"github.com/maxperrimond/kurin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type (
	Adapter struct {
		srv     *http.Server
		port    string
		healthy bool
	}
)

func NewHTTPAdapter(handler http.Handler, port string) kurin.Adapter {
	adapter := &Adapter{
		port:    port,
		healthy: true,
	}

	inFlight := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "in_flight_requests",
		Help: "A gauge of requests currently being served by the wrapped handler.",
	})
	totalCount := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "A counter for requests to the wrapped handler.",
		},
		[]string{"code", "method"},
	)
	durationHist := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "response_duration_seconds",
			Help:        "A histogram of request latencies.",
			Buckets:     prometheus.DefBuckets,
			ConstLabels: prometheus.Labels{"handler": "api"},
		},
		[]string{"method"},
	)
	prometheus.MustRegister(inFlight, totalCount, durationHist)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if adapter.healthy {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/",
		promhttp.InstrumentHandlerInFlight(inFlight,
			promhttp.InstrumentHandlerCounter(totalCount,
				promhttp.InstrumentHandlerDuration(durationHist, handler))))

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
	err := adapter.srv.Close()
	if err != nil {
		log.Println(err)
	}
}

func (adapter *Adapter) Healthy() bool {
	return adapter.healthy
}

func (adapter *Adapter) ListenFailure(ce <-chan error) {
	go func() {
		err := <-ce
		if err != nil {
			adapter.healthy = false
		}
	}()
}
