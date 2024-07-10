package httpx

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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

type ErrorFunc func(http.ResponseWriter, *http.Request, error)

func onError(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), 500)
	log.Panicf("%v: %v\n", r.URL.Path, err)
}

type HandlerFunc func(http.ResponseWriter, *http.Request) error

type MiddlewareFunc func(http.Handler) http.Handler

func wrap(h http.Handler, middlewares ...MiddlewareFunc) http.Handler {
	for _, mid := range middlewares {
		h = mid(h)
	}
	return h
}

type handler struct {
	onFunc  HandlerFunc
	onError ErrorFunc
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.onFunc(w, r)
	if err != nil {
		h.onError(w, r, err)
	}
}

type ServeMux struct {
	http.ServeMux
	onError     ErrorFunc
	middlewares []MiddlewareFunc
}

func NewServeMux(middlewares ...MiddlewareFunc) *ServeMux {
	return &ServeMux{
		ServeMux:    *http.NewServeMux(),
		onError:     onError,
		middlewares: middlewares,
	}
}

func (mux *ServeMux) Add(middlewares ...MiddlewareFunc) {
	mux.middlewares = append(mux.middlewares, middlewares...)
}

func (mux *ServeMux) OnError(onError ErrorFunc) {
	mux.onError = onError
}

func (mux *ServeMux) With(middlewares ...MiddlewareFunc) http.Handler {
	var httph http.Handler = mux
	for _, mid := range middlewares {
		httph = mid(httph)
	}
	return httph
}

func (mux *ServeMux) Handle(pattern string, handler http.Handler, middlewares ...MiddlewareFunc) {
	mux.ServeMux.Handle(pattern, wrap(wrap(handler, mux.middlewares...), middlewares...))

}
func (mux *ServeMux) HandleFunc(pattern string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	mux.ServeMux.Handle(pattern, wrap(wrap(handler, mux.middlewares...), middlewares...))
}
func (mux *ServeMux) Handlex(pattern string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	mux.ServeMux.Handle(pattern, wrap(wrap(handler{onFunc: handlerFunc, onError: mux.onError}, mux.middlewares...), middlewares...))
}

func (mux *ServeMux) NewGroup(path string, middlewares ...MiddlewareFunc) *Group {
	return &Group{
		mux:         mux,
		path:        path,
		middlewares: middlewares,
		onError:     onError,
	}
}

type Group struct {
	mux         *ServeMux
	path        string
	middlewares []MiddlewareFunc
	onError     ErrorFunc
}

func (g *Group) Add(middlewares ...MiddlewareFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *Group) OnError(onError ErrorFunc) {
	g.onError = onError
}

func (g *Group) Handle(pattern string, handler http.Handler, middlewares ...MiddlewareFunc) {
	g.mux.Handle(pattern, wrap(wrap(handler, g.middlewares...), middlewares...))

}
func (g *Group) HandleFunc(pattern string, handler http.HandlerFunc, middlewares ...MiddlewareFunc) {
	g.mux.Handle(pattern, wrap(wrap(handler, g.middlewares...), middlewares...))
}

func (g *Group) Handlex(pattern string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	method := ""
	split := strings.Split(pattern, " ")
	if len(split) == 2 {
		method = split[0]
		pattern = split[1]
	}
	g.mux.Handlex(fmt.Sprintf("%v %v%v", method, g.path, pattern), handler, append(g.middlewares, middlewares...)...)
}
