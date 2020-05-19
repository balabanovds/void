package apiserver

import (
	"github.com/balabanovds/void/internal/domain"
	"github.com/balabanovds/void/internal/server"
	"github.com/rs/zerolog"
)

// APIServer RESTful implementation
type APIServer struct {
	config  *server.Config
	logger  *zerolog.Logger
	storage domain.Storage
}

// New API server instance
func New(config *server.Config, storage domain.Storage, logger *zerolog.Logger) *APIServer {
	return &APIServer{
		config:  config,
		storage: storage,
		logger:  logger,
	}
}

// Start API server instance
func (s *APIServer) Start() {}
