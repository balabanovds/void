package apiserver

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *APIServer) routes() http.Handler {
	mux := mux.NewRouter()

	mux.HandleFunc("/", s.handleHome())

	return mux
}
