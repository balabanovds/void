package apiserver

import "net/http"

func (s *APIServer) handleSKUD() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("SKUD hit"))
	}
}
