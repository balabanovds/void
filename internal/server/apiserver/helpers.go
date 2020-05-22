package apiserver

import (
	"encoding/json"
	"github.com/balabanovds/void/internal/server/ctxHelper"
	"net/http"
)

func (s *APIServer) clientError(w http.ResponseWriter, r *http.Request, code int, err error) {
	email, ok := r.Context().Value(ctxHelper.CtxKeyEmail).(string)
	if !ok {
		email = "unknown"
	}
	s.log.Warn().
		Err(err).
		Msgf("code: %d; user: %s; from: %s; %s %s",
			code, email, r.RemoteAddr, r.Method, r.RequestURI)
	s.respond(w, code, map[string]string{"error": err.Error()})
}

func (s *APIServer) serverError(w http.ResponseWriter, err error) {
	s.log.Error().Caller(1).Err(err).Send()
	s.respond(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}

func (s *APIServer) respond(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
