package apiserver

import "net/http"

func (s *APIServer) handleWiki() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		renderStub(w, r)
	}
}
