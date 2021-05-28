package main

import (
	"crypto/tls"
	"os"
)

type Config struct {
	DBPath         string
	DBDriver       string
	LogLevel       string
	SMTPServerAddr string
	POP3ServerAddr string
	Username       string
	Password       string
	TLSConfig      *tls.Config
	BindAddr       string
}

func NewConfig() *Config {
	return &Config{
		DBPath:         "./mails.db",
		DBDriver:       "sqlite3",
		LogLevel:       "debug",
		SMTPServerAddr: "smtp.mailtrap.io:25",
		POP3ServerAddr: "pop3.mailtrap.io:1100",
		TLSConfig:      &tls.Config{ServerName: "mailtrap.io"},
		Username:       os.Getenv("MAILTRAP_USERNAME"),
		Password:       os.Getenv("MAILTRAP_PASSWORD"),
		BindAddr:       ":8080",
	}
}
