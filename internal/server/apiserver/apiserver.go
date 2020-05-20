package apiserver

import (
	"fmt"
	"net/http"

	"github.com/balabanovds/void/internal/domain"
	"github.com/balabanovds/void/internal/server"
	"github.com/balabanovds/void/internal/service"
	"github.com/rs/zerolog"
)

// APIServer RESTful implementation
type APIServer struct {
	config  *server.Config
	service *service.Service
	log     zerolog.Logger
	debug   zerolog.Logger
}

// New API server instance
func New(config *server.Config, storage domain.Storage, logger zerolog.Logger) *APIServer {
	l := logger.With().Str("server", "API").Logger()
	return &APIServer{
		config:  config,
		service: service.New(storage, logger),
		log:     l,
		debug:   l.With().Caller().Stack().Logger(),
	}
}

// Start API server instance
func (s *APIServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Hostname, s.config.Port)

	s.log.Info().Msgf("API server is running on http://%s", addr)
	return http.ListenAndServe(addr, s.routes())
}
