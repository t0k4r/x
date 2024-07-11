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

func (h Handler) Wrap(middlewares ...MiddlewareFunc) http.Handler {
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

func (mux *ServeMux) Add(middlewares ...MiddlewareFunc) {
	mux.middlewares = append(mux.middlewares, middlewares...)
}

func (mux *ServeMux) Wrap(middlewares ...MiddlewareFunc) http.Handler {
	return Wrap(mux, middlewares...)
}

func (mux *ServeMux) Handle(pattern string, handler http.Handler, middlewares ...MiddlewareFunc) {
	mux.ServeMux.Handle(pattern, Wrap(Wrap(handler, mux.middlewares...), middlewares...))

}
func (mux *ServeMux) HandleFunc(pattern string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	mux.Handle(pattern, handler, middlewares...)
}
func (mux *ServeMux) Handlex(pattern string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	mux.Handle(pattern, Handler{OnFunc: handlerFunc, OnError: mux.OnError}, middlewares...)
}

func (mux *ServeMux) With(routes func(*WithHandler), middlewares ...MiddlewareFunc) {
	routes(&WithHandler{
		mux:         mux,
		middlewares: middlewares,
		OnError:     mux.OnError,
	})
}

type WithHandler struct {
	mux         *ServeMux
	middlewares []MiddlewareFunc
	OnError     ErrorFunc
}

func (wh *WithHandler) Add(middlewares ...MiddlewareFunc) {
	wh.middlewares = append(wh.middlewares, middlewares...)
}

func (wh *WithHandler) Handle(pattern string, handler http.Handler, middlewares ...MiddlewareFunc) {
	wh.mux.Handle(pattern, Wrap(handler, wh.middlewares...), middlewares...)
}

func (wh *WithHandler) HandleFunc(pattern string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	wh.Handle(pattern, handler, middlewares...)
}

func (wh *WithHandler) Handlex(pattern string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	wh.Handle(pattern, Handler{OnFunc: handler, OnError: wh.OnError}, middlewares...)
}

func (wh *WithHandler) With(routes func(*WithHandler), middlewares ...MiddlewareFunc) {
	routes(&WithHandler{
		mux:         wh.mux,
		middlewares: wh.middlewares,
		OnError:     wh.OnError,
	})
}
