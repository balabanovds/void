package apiserver

import (
	"github.com/gorilla/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *APIServer) routes() http.Handler {
	r := mux.NewRouter()

	r.Use(jsonContent)
	r.Use(s.setRequestID)
	r.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	r.HandleFunc("/", s.handleHome())
	r.HandleFunc("/signup", s.handleSignUp()).Methods("POST")
	r.HandleFunc("/login", s.handleLogin()).Methods("POST")

	r.HandleFunc("/wiki", s.handleWiki()).Methods("GET")

	private := r.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/skud2", s.handleSKUD()).Methods("POST")

	return r
}
