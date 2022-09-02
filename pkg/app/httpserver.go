package app

import "net/http"

type Server interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	ListenAndServe(addr string, handler http.Handler) error
}

type HttpServer struct{}

func (HttpServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, handler)
}

func (HttpServer) ListenAndServe(addr string, handler http.Handler) error {
	return http.ListenAndServe(addr, handler)
}

// MockServer is implementing Server interface just for testing purpose
type MockServer struct{}

func (MockServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
}

func (MockServer) ListenAndServe(addr string, handler http.Handler) error {
	return nil
}
