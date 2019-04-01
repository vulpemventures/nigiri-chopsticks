package middleware

import (
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Logger is a middleware that logs every request/response
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		res := newResponseWriter(w)

		next.ServeHTTP(res, r)

		log.WithFields(log.Fields{
			"start_time": start.Format(time.RFC3339),
			"duration":   time.Since(start),
			"method":     r.Method,
			"hostname":   r.Host,
			"path":       r.URL.Path,
			"status":     res.Status(),
			"response":   prettyJSON(res.Body()),
		}).Info("New request from remote address: ", r.RemoteAddr)
	})
}

func prettyJSON(b []byte) string {
	replacer := strings.NewReplacer("\"", "'", "\n", "")
	return replacer.Replace(string(b))
}
