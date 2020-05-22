package apiserver

import (
	"github.com/gorilla/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *APIServer) routes() http.Handler {
	router := mux.NewRouter()

	router.Use(jsonContent)
	router.Use(s.setRequestID)
	router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	router.HandleFunc("/", s.handleHome())
	router.HandleFunc("/signup", s.handleSignUp()).Methods("POST")

	return router
}
