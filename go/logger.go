package swagger

import (
	"net/http"
	"time"
)

// Logger for HTTP calls
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		lg.Debugf("%s %s %s %s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start))
	})
}
