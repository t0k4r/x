package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func NoContent(w http.ResponseWriter, code int) error {
	w.WriteHeader(code)
	return nil
}

func Error(w http.ResponseWriter, err error, code int) error {
	http.Error(w, err.Error(), code)
	return nil
}

func Text(w http.ResponseWriter, text string, code int) error {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "text/plain")
	_, err := fmt.Fprint(w, text)
	return err

}

func Html(w http.ResponseWriter, html string, code int) error {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "text/html")
	_, err := fmt.Fprint(w, html)
	return err

}

func Json(w http.ResponseWriter, v any, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

type MiddlewareFunc func(http.Handler) http.Handler

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := f(w, r)
	if err != nil {
		Error(w, err, 500)
	}
}

func wrap(handler http.Handler, middlewares ...MiddlewareFunc) http.Handler {
	for _, mid := range middlewares {
		handler = mid(handler)
	}
	return handler
}

type ServeMux struct {
	http.ServeMux
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		ServeMux: *http.NewServeMux(),
	}
}

func (mux *ServeMux) With(middlewares ...MiddlewareFunc) http.Handler {
	return wrap(mux, middlewares...)
}

func (mux *ServeMux) Handle(pattern string, handler http.Handler, middlewares ...MiddlewareFunc) {
	mux.ServeMux.Handle(pattern, wrap(handler, middlewares...))
}
func (mux *ServeMux) HandleFunc(pattern string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	mux.ServeMux.Handle(pattern, wrap(handler, middlewares...))
}

func (mux *ServeMux) Handlex(pattern string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	mux.ServeMux.Handle(pattern, wrap(handler, middlewares...))
}
