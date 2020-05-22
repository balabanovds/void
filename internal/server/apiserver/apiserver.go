package apiserver

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"net/http"

	"github.com/balabanovds/void/internal/domain"
	"github.com/balabanovds/void/internal/server"
	"github.com/balabanovds/void/internal/service"
	"github.com/rs/zerolog"
)

// APIServer RESTful implementation
type APIServer struct {
	config       *server.Config
	service      *service.Service
	sessionStore sessions.Store
	log          zerolog.Logger
}

var (
	sessionKey   string
	sessionEmail = "user_email"
)

// New API server instance
func New(config *server.Config, storage domain.Storage, logger zerolog.Logger) *APIServer {
	sessionKey = uuid.New().String()
	return &APIServer{
		config:       config,
		service:      service.New(storage, logger),
		sessionStore: sessions.NewCookieStore([]byte(sessionKey)),
		log:          logger.With().Str("server", "api").Logger(),
	}
}

// Start API server instance
func (s *APIServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Hostname, s.config.Port)

	s.log.Info().Msgf("API server is running on http://%s", addr)
	return http.ListenAndServe(addr, s.routes())
}
