package httpx

import (
	"encoding/json"
	"fmt"
	"log"
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

var DefalutOnError ErrorFunc = onError

type ErrorFunc func(http.ResponseWriter, *http.Request, error)

func onError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), 500)
	log.Panicf("%v: %v\n", r.URL.Path, err)
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

type MiddlewareFunc func(http.Handler) http.Handler

func Wrap(h http.Handler, middlewares ...MiddlewareFunc) http.Handler {
	for _, mid := range middlewares {
		h = mid(h)
	}
	return h
}

type Handler struct {
	OnFunc  HandlerFunc
	OnError ErrorFunc
}

func (h Handler) With(middlewares ...MiddlewareFunc) http.Handler {
	return Wrap(h, middlewares...)
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.OnFunc(w, r)
	if err != nil {
		h.OnError(w, r, err)
	}
}

type ServeMux struct {
	http.ServeMux
	OnError     ErrorFunc
	middlewares []MiddlewareFunc
}

func NewServeMux(middlewares ...MiddlewareFunc) *ServeMux {
	return &ServeMux{
		ServeMux:    *http.NewServeMux(),
		OnError:     DefalutOnError,
		middlewares: middlewares,
	}
}

// push middleware to stack this middlewares are only applied to routes added below using Handle* methods
func (mux *ServeMux) Push(middlewares ...MiddlewareFunc) {
	mux.middlewares = append(mux.middlewares, middlewares...)
}

// remove last pushed midleware
func (mux *ServeMux) Pop() {
	mux.middlewares = mux.middlewares[:len(mux.middlewares)-1]
}
func (mux *ServeMux) With(middlewares ...MiddlewareFunc) http.Handler {
	return Wrap(mux, middlewares...)
}

func (mux *ServeMux) Handle(pattern string, handler http.Handler, middlewares ...MiddlewareFunc) {
	mux.ServeMux.Handle(pattern, Wrap(Wrap(handler, mux.middlewares...), middlewares...))

}
func (mux *ServeMux) HandleFunc(pattern string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	mux.ServeMux.Handle(pattern, Wrap(Wrap(handler, mux.middlewares...), middlewares...))
}
func (mux *ServeMux) Handlex(pattern string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	mux.ServeMux.Handle(pattern, Wrap(Wrap(Handler{OnFunc: handlerFunc, OnError: mux.OnError}, mux.middlewares...), middlewares...))
}
