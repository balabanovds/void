package apiserver

import (
	"encoding/json"
	"errors"
	"github.com/balabanovds/void/internal/domain"
	"github.com/balabanovds/void/internal/models"
	"github.com/balabanovds/void/internal/service"
	"net/http"
)

func (s *APIServer) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("home hit"))
	}
}

func (s *APIServer) handleSignUp() http.HandlerFunc {
	type request struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"password_confirm"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.clientError(w, r, http.StatusBadRequest, err)
			return
		}

		user := models.User{Email: req.Email, Password: req.Password}
		if err := user.Validate(); err != nil {
			s.clientError(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		user, err := s.service.Users().Create(req.Email, req.Password, req.PasswordConfirm)
		if err != nil {
			if errors.Is(err, service.ErrPasswdNotMatch) || errors.Is(err, domain.ErrDuplicateEmail) {
				s.clientError(w, r, http.StatusConflict, err)
				return
			}
			s.serverError(w, err)
			return
		}
		s.log.Info().Msgf("new user %s created", req.Email)
		s.respond(w, http.StatusCreated, &user)
	}
}
