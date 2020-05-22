package apiserver

import (
	"fmt"
	"github.com/gorilla/securecookie"
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
	sessionKey []byte
)

// New API server instance
func New(config *server.Config, storage domain.Storage, logger zerolog.Logger) *APIServer {
	sessionKey = securecookie.GenerateRandomKey(32)
	return &APIServer{
		config:       config,
		service:      service.New(storage, logger),
		sessionStore: sessions.NewCookieStore(sessionKey),
		log:          logger.With().Str("server", "api").Logger(),
	}
}

// Start API server instance
func (s *APIServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Hostname, s.config.Port)

	s.log.Info().Msgf("API server is running on http://%s", addr)
	return http.ListenAndServe(addr, s.routes())
}
