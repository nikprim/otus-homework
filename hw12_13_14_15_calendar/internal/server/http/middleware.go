package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func loggingMiddleware() mux.MiddlewareFunc {
	return hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		l := log.With().Fields(map[string]interface{}{
			"remote_addr": r.RemoteAddr,
			"method":      r.Method,
			"protocol":    r.Proto,
			"url":         r.URL.String(),
			"user_agent":  r.Header.Get("User-Agent"),
			"referer":     r.Header.Get("Referer"),
			"status":      status,
			"size":        fmt.Sprintf("%d bytes", size),
			"duration":    fmt.Sprintf("%d ms", duration.Milliseconds()),
		}).Logger()

		l.Info().Send()
	})
}
