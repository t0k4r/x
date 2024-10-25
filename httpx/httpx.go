package httpx

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"sync/atomic"
	"time"
)

// log HandlerFunc error as Debug Logger middleware as Info and REcoverer as Error
var defaultLogger atomic.Pointer[slog.Logger]

func init() {
	defaultLogger.Store(slog.Default())
}
func SetLogger(l *slog.Logger) {
	defaultLogger.Store(l)
}

func Json(w http.ResponseWriter, v any, code int) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func Plain(w http.ResponseWriter, text string, code int) error {
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(code)
	_, err := fmt.Fprint(w, text)
	return err
}

func NoContent(w http.ResponseWriter, code int) error {
	w.WriteHeader(code)
	return nil
}

func Error(w http.ResponseWriter, err error, code int) error {
	http.Error(w, err.Error(), code)
	return err
}

type funcResponseWriter struct {
	http.ResponseWriter
	statusCode int
	header     bool
}

func (rw *funcResponseWriter) WriteHeader(statusCode int) {
	rw.WriteHeader(statusCode)
	rw.statusCode = statusCode
	rw.header = true
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func (hf HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := &funcResponseWriter{ResponseWriter: w}
	err := hf(rw, r)
	if err != nil {
		if !rw.header {
			w.WriteHeader(http.StatusInternalServerError)
		}
		defaultLogger.Load().DebugContext(r.Context(), r.URL.Path, "code", rw.statusCode, "error", err.Error())
	}
}

type MiddlewareFunc func(http.Handler) http.Handler

func Wrap(handler http.Handler, middlewares ...MiddlewareFunc) http.Handler {
	for _, mid := range middlewares {
		handler = mid(handler)
	}
	return handler
}

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				defaultLogger.Load().ErrorContext(r.Context(), r.URL.Path, "panic", fmt.Sprint(rec), "stack", debug.Stack())
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type logResponseWriter struct {
	http.ResponseWriter
	now        time.Time
	statusCode int
}

func (rw *logResponseWriter) WriteHeader(statusCode int) {
	rw.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &logResponseWriter{ResponseWriter: w, now: time.Now(), statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)
		defaultLogger.Load().InfoContext(r.Context(), r.URL.Path, "code", rw.statusCode, "time", time.Since(rw.now))
	})

}

type ServeMux struct {
	*http.ServeMux
}

func NewServeMux() ServeMux {
	return ServeMux{http.NewServeMux()}
}
func (mux ServeMux) HandleFunc(pattern string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	mux.Handle(pattern, Wrap(handler, middlewares...))
}
