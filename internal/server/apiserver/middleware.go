package apiserver

import (
	"context"
	"github.com/balabanovds/void/internal/server/ctxHelper"
	"github.com/balabanovds/void/internal/service"
	"github.com/google/uuid"
	"net/http"
)

func jsonContent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		ctx := context.WithValue(r.Context(), ctxHelper.CtxKeyRequestID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *APIServer) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func (s *APIServer) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionKey)
		if err != nil {
			s.serverError(w, err)
			return
		}

		email, ok := session.Values[sessionEmail]
		if !ok {
			s.clientError(w, r, http.StatusUnauthorized, service.ErrNotAllowed)
			return
		}

		emailStr, ok := email.(string)
		if !ok {
			s.serverError(w, err)
			return
		}

		_, err = s.service.Users().GetByEmail(emailStr)
		if err != nil {
			s.clientError(w, r, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), ctxHelper.CtxKeyEmail, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
