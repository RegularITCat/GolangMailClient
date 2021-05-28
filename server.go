package main

import (
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

type Server struct {
	config *Config
	router *mux.Router
}

func NewServer() *Server {
	return &Server{
		config: NewConfig(),
		router: mux.NewRouter(),
	}
}

func (s *Server) Start() error {
	s.configureRouter()
	return nil
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/hello", s.handleHello())
}

func (s *Server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "Hello, World!")
		if err != nil {
			log.Fatal(err)
		}
	}
}
