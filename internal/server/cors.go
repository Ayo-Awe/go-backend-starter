package server

import (
	"net/http"
	"strings"

	"github.com/go-chi/cors"
)

func (s Server) CorsMiddleware() func(next http.Handler) http.Handler {
	if s.corsAllowedOrigins == nil {
		s.logger.Info("CORS – all origins are allowed")
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	allowedMethods := []string{
		http.MethodHead,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}

	s.logger.Info("CORS setup",
		"allowedOrigins", strings.Join(s.corsAllowedOrigins, ","),
		"headers", strings.Join(allowedMethods, ","),
	)

	return cors.New(cors.Options{
		AllowedOrigins: s.corsAllowedOrigins,
		AllowedMethods: allowedMethods}).Handler
}
