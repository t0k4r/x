package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

type LogRespnseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (l *LogRespnseWriter) WriteHeader(statusCode int) {
	l.statusCode = statusCode
	l.ResponseWriter.WriteHeader(statusCode)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		lrw := &LogRespnseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(lrw, r)
		slog.Info(http.StatusText(lrw.statusCode), "path", r.URL.Path, "code", lrw.statusCode, "time", time.Since(now))
	})
}

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				slog.Error(fmt.Sprintf("recovered panic %v", rec), "path", r.URL.Path, "stack", debug.Stack())
			}
		}()
		next.ServeHTTP(w, r)
	})
}
