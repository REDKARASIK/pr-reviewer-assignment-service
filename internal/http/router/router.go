package router

import "net/http"

type Router struct {
	mux  *http.ServeMux
	base string
}

func NewRouter() *Router {
	return &Router{mux: http.NewServeMux()}
}

func (r *Router) Handler() http.Handler {
	return r.mux
}

func (r *Router) Group(prefix string) *Router {
	return &Router{
		mux:  r.mux,
		base: r.base + prefix,
	}
}

func (r *Router) Handle(pattern string, h http.Handler) {
	r.mux.Handle(pattern, h)
}

func (r *Router) HandleFunc(pattern string, h func(http.ResponseWriter, *http.Request)) {
	r.mux.HandleFunc(pattern, h)
}

func (r *Router) GET(path string, h http.HandlerFunc) {
	r.handle(http.MethodGet, path, h)
}
func (r *Router) POST(path string, h http.HandlerFunc) {
	r.handle(http.MethodPost, path, h)
}
func (r *Router) PUT(path string, h http.HandlerFunc) {
	r.handle(http.MethodPut, path, h)
}
func (r *Router) DELETE(path string, h http.HandlerFunc) {
	r.handle(http.MethodDelete, path, h)
}

func (r *Router) fullPath(path string) string {
	if path == "" {
		return r.base
	}
	return r.base + path
}

func (r *Router) handle(method, path string, h http.HandlerFunc) {
	pattern := method + " " + r.fullPath(path)
	r.mux.Handle(pattern, h)
}
