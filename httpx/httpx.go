package httpx

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func LogErr(r *http.Request, err error) {
	slog.Error(err.Error(), "path", r.URL.Path)
}

func Json(w http.ResponseWriter, v any, code int) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func Text(w http.ResponseWriter, text string, code int) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	fmt.Fprint(w, text)
	return nil
}

func Html(w http.ResponseWriter, html string, code int) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	fmt.Fprint(w, html)
	return nil
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
		LogErr(r, err)
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
