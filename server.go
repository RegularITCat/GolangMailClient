package main

import (
	"encoding/json"
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
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/hello", s.handleHello())
	s.router.HandleFunc("/api/mails", s.handleAPIGetAllMails()).Methods(http.MethodGet)
}

func (s *Server) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "Hello, World!")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *Server) handleAPIGetAllMails() http.HandlerFunc {
	mm := MailMap{
		DBPath:   s.config.DBPath,
		DBDriver: s.config.DBDriver,
	}
	mails, err := mm.SelectAll()
	if err != nil {
		log.Fatal(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(mails)
		if err != nil {
			log.Fatal(err)
		}

	}
}
