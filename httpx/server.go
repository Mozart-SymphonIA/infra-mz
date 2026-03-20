package httpx

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Endpoint interface {
	Pattern() string
	Handler() http.Handler
}

type HealthChecker interface {
	Health(ctx context.Context) error
}

type HealthyService struct {
	Name    string
	Service HealthChecker
}

type endpoint struct {
	pattern string
	handler http.Handler
}

var _ Endpoint = (*endpoint)(nil)

func (e endpoint) Pattern() string       { return e.pattern }
func (e endpoint) Handler() http.Handler { return e.handler }

func NewEndpoint(pattern string, h http.Handler) Endpoint {
	return endpoint{
		pattern: pattern,
		handler: h,
	}
}

func NewEndpointFunc(pattern string, fn func(http.ResponseWriter, *http.Request)) Endpoint {
	return endpoint{
		pattern: pattern,
		handler: http.HandlerFunc(fn),
	}
}

func NewHTTPServer(port string, healthyServices []HealthyService, endpoints ...Endpoint) *http.Server {

	mux := http.NewServeMux()
	for _, ep := range endpoints {
		mux.Handle(ep.Pattern(), ep.Handler())
	}

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		for _, s := range healthyServices {
			if err := s.Service.Health(r.Context()); err != nil {
				http.Error(w, fmt.Sprintf("%s not ready: %v", s.Name, err), http.StatusServiceUnavailable)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	})

	return &http.Server{
		Addr:              port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}
