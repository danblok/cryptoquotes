package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/danblok/cryptoquotes/pkg/config"
	"github.com/danblok/cryptoquotes/pkg/types"
)

// HTTPServer implementation for QuotesSource.
type HTTPServer struct {
	svc    types.QuotesSource
	srv    *http.Server
	config *config.HTTPConfig
}

// HTTPErrResponse is a response body.
type HTTPErrResponse struct {
	Error any `json:"error"`
}

// HTTPHandlerFunc is a helper handler func.
type HTTPHandlerFunc func(context.Context, http.ResponseWriter, *http.Request) error

// NewHTTPServer constructs new HTTPServer that signs and validates tokens via HTTP.
func NewHTTPServer(svc types.QuotesSource, config *config.HTTPConfig) *HTTPServer {
	return &HTTPServer{
		config: config,
		svc:    svc,
		srv: &http.Server{
			Addr:        fmt.Sprintf(":%d", config.Port()),
			ReadTimeout: 5 * time.Second,
			IdleTimeout: 5 * time.Second,
		},
	}
}

// ListenAndServe runs the server.
func (s *HTTPServer) ListenAndServe() error {
	mux := http.NewServeMux()
	mux.Handle("GET /", makeHTTPHandler(s.handleQuotes))
	s.srv.Handler = mux
	return s.srv.ListenAndServe()
}

// ListenAndServeTLS runs the server with TLS enabled.
func (s *HTTPServer) ListenAndServeTLS() error {
	mux := http.NewServeMux()
	mux.Handle("GET /", makeHTTPHandler(s.handleQuotes))
	s.srv.Handler = mux
	s.srv.TLSConfig = s.config.TLSConfig()
	return s.srv.ListenAndServeTLS("", "")
}

// Attaches request_id to the context and returns http.Handler.
func makeHTTPHandler(fn HTTPHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), types.RequestIDKey, uuid.NewString())
		ctx = context.WithValue(ctx, types.TransportTypeKey, types.HTTPTransport)
		ctx = context.WithValue(ctx, types.RemoteAddrKey, r.RemoteAddr)

		if err := fn(ctx, w, r); err != nil {
			_ = writeJSON(w, http.StatusBadRequest, HTTPErrResponse{Error: err.Error()})
		}
	}
}

// handleQuotes handles the '/' path.
func (s *HTTPServer) handleQuotes(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query().Get("symbol")
	names := strings.Split(query, ",")

	quotes, err := s.svc.Quotes(ctx, names...)
	if err != nil {
		return err
	}

	err = writeJSON(w, http.StatusOK, quotes)
	if err != nil {
		return err
	}

	return nil
}

// Helper func for responding with JSON.
func writeJSON(w http.ResponseWriter, code int, body any) error {
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(body)
}
