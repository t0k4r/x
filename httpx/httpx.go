package httpx

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func ErrorLog(w http.ResponseWriter, r *http.Request, err error, msg string, code int) error {
	slog.Error(err.Error(), "path", r.URL.Path)
	http.Error(w, msg, code)
	return nil
}

func Json(w http.ResponseWriter, v any, code int) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func Text(w http.ResponseWriter, text string, code int) error {
	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(code)
	_, err := fmt.Fprint(w, text)
	return err
}

func Html(w http.ResponseWriter, html string, code int) error {
	w.Header().Set("content-type", "text/html")
	w.WriteHeader(code)
	_, err := fmt.Fprint(w, html)
	return err
}

func Empty(w http.ResponseWriter, code int) error {
	w.WriteHeader(code)
	return nil
}

func Error(w http.ResponseWriter, err error, code int) error {
	http.Error(w, err.Error(), code)
	return nil
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func (hf HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := hf(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error(err.Error(), "path", r.URL.Path)
	}
}

type MiddlewareFunc func(http.Handler) http.Handler

func wrap(h http.Handler, middlewares ...MiddlewareFunc) http.Handler {
	for _, mid := range middlewares {
		h = mid(h)
	}
	return h
}

type ServeMux struct {
	http.ServeMux
}

func NewServeMux() *ServeMux {
	return &ServeMux{ServeMux: *http.NewServeMux()}
}

func (mux *ServeMux) Handle(pattern string, handler http.Handler, middlewares ...MiddlewareFunc) {
	mux.ServeMux.Handle(pattern, wrap(handler, middlewares...))
}
func (mux *ServeMux) HandleFunc(pattern string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	mux.Handle(pattern, handler, middlewares...)
}
func (mux *ServeMux) Handlex(pattern string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	mux.Handle(pattern, handler, middlewares...)
}

func ListenAndServe(addr string, handler http.Handler, middlewares ...MiddlewareFunc) error {
	return http.ListenAndServe(addr, wrap(handler, middlewares...))
}
