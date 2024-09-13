package middleware

import (
	"fmt"
	"net/http"
	"time"

	"soul-connection.com/api/src/lib"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *wrappedWriter) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		lib.ServerLog("INFO", fmt.Sprintf("%s [%d] %s %s", r.Method, wrapped.statusCode, r.URL.Path, time.Since(start)))
	})
}
