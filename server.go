package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/simia-tech/go-pop3"
	"io"
	"log"
	"net"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"
)

type Server struct {
	config  *Config
	router  *mux.Router
	mailMap MailMap
}

func NewServer() *Server {
	server := &Server{
		config:  NewConfig(),
		router:  mux.NewRouter(),
		mailMap: MailMap{},
	}
	server.mailMap.DBPath = server.config.DBPath
	server.mailMap.DBDriver = server.config.DBDriver
	return server
}

func (s *Server) Start() error {
	//TODO Check database exist and if not create it!
	if _, err := os.Stat(s.config.DBPath); err != nil {
		err = CreateDatabase(s.config)
		if err != nil {
			log.Fatal(err)
		}
		err = CreateTable(s.config)
		if err != nil {
			log.Fatal(err)
		}
	}
	s.configureRouter()
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *Server) configureRouter() {
	//TODO Create all API endpoints to database function's
	//TODO Make hard fetch from mailserver
	s.router.HandleFunc("/hello", s.handleHello())
	s.router.HandleFunc("/api/mail", s.handleAPIGetAllMails()).Methods(http.MethodGet)
	s.router.HandleFunc("/api/mail/sync", s.handleSyncMailbox()).Methods(http.MethodGet)
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

	return func(w http.ResponseWriter, r *http.Request) {
		mails, err := s.mailMap.SelectAll()
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(mails)
		if err != nil {
			log.Fatal(err)
		}

	}
}

func (s *Server) handleSyncMailbox() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		pop3Conn, err := net.Dial("tcp", s.config.POP3ServerAddr)
		if err != nil {
			log.Fatal(err)
		}

		c, err := pop3.NewClient(pop3Conn, pop3.UseTLS(s.config.TLSConfig), pop3.UseTimeout(time.Second*1))
		if err != nil {
			log.Fatal(err)
		}

		err = c.Auth(s.config.Username, s.config.Password)
		if err != nil {
			log.Fatal(err)
		}

		data, err := c.ListAll()
		if err != nil {
			log.Fatal(err)
		}

		messages := make(map[int]string)
		for _, message := range data {
			str, err := c.Retr(message.Seq)
			if err != nil {
				log.Fatal(err)
			}
			messages[int(message.Seq)] = str
		}

		err = c.Quit()
		if err != nil {
			log.Fatal(err)
		}

		oldMails, err := s.mailMap.SelectAll()
		if err != nil {
			log.Fatal(err)
		}

		for key, value := range messages {
			if val, ok := oldMails[key]; ok {
				reader := strings.NewReader(value)
				msg, err := mail.ReadMessage(reader)
				if err != nil {
					log.Fatal(err)
				}
				header := msg.Header
				val.To = header.Get("To")
				val.From = header.Get("From")
				val.Subject = header.Get("Subject")
				val.FullText = value
				err = s.mailMap.Update(val)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				reader := strings.NewReader(value)
				msg, err := mail.ReadMessage(reader)
				if err != nil {
					log.Fatal(err)
				}
				header := msg.Header
				m := &Mail{
					From:     header.Get("From"),
					To:       header.Get("To"),
					Subject:  header.Get("Subject"),
					FullText: value,
				}
				err = s.mailMap.Insert(m)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]string{
			"status": "done",
		})
		if err != nil {
			log.Fatal(err)
		}

	}
}
