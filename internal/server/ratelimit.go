package server

import (
	"log/slog"
	"net"
	"net/http"
	"strings"

	"github.com/ayo-awe/go-backend-starter/internal/oapi"
	"github.com/go-chi/httprate"
	"github.com/tomasen/realip"
)

func (s Server) RateLimitMiddleware() func(next http.Handler) http.Handler {
	if !s.rateLimitEnabled {
		s.logger.Info("rate limit disabled")
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	s.logger.Info("rate limit enabled",
		slog.Int("rate", s.rateLimitRate),
		slog.Duration("windowLength", s.rateLimitWindow),
		slog.String("whitelistedIPs", strings.Join(s.rateLimitWhitelist, ",")),
	)

	// Create a map for faster IP lookup
	whitelistedIPMap := make(map[string]bool)
	for _, ip := range s.rateLimitWhitelist {
		whitelistedIPMap[ip] = true
	}

	return httprate.Limit(
		s.rateLimitRate,
		s.rateLimitWindow,
		httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
			ip := realip.FromRequest(r)

			// Check if the IP is in the whitelist
			if whitelistedIPMap[ip] {
				return "whitelisted", nil
			}

			// Check if the IP is in a whitelisted range
			for whitelistedIP := range whitelistedIPMap {
				if strings.Contains(whitelistedIP, "/") {
					_, ipNet, err := net.ParseCIDR(whitelistedIP)
					if err == nil && ipNet.Contains(net.ParseIP(ip)) {
						return "whitelisted", nil
					}
				}
			}

			return ip, nil
		}),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			s.errorResponse(w, r,
				oapi.ErrorCodeRateLimitExceeded,
				http.StatusTooManyRequests,
				"too many requests, please slow down",
			)
		}),
	)
}
