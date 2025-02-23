package server

import "net/http"

func (s Server) Welcome(w http.ResponseWriter, r *http.Request) {
	s.sendJSONResponse(w, r, http.StatusNoContent, nil)
}
