package main

import (
	"crypto/tls"
	"fmt"
	"github.com/RegularITCat/GolangMailClient/pop3"
	"log"
	"net"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

/*type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}


func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}


func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}*/

func main() {
	host := "smtp.mailtrap.io"
	port := 25
	to := "to@example.com"
	msg := []byte("To: to@example.com\r\n" +
		"Subject: Why are you not using Mailtrap yet?\r\n" +
		"\r\n" +
		"Here’s the space for our great sales pitch\r\n")
	auth := smtp.PlainAuth("", os.Getenv("MAILTRAP_USERNAME"), os.Getenv("MAILTRAP_PASSWORD"), host)
	conf := &tls.Config{ServerName: host}

	conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}

	smtpCon, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Fatal(err)
	}

	err = smtpCon.StartTLS(conf)
	if err != nil {
		log.Fatal(err)
	}

	err = smtpCon.Auth(auth)
	if err != nil {
		log.Fatal(err)
	}

	err = smtpCon.Mail("example@example.com")
	if err != nil {
		log.Fatal(err)
	}

	err = smtpCon.Rcpt(to)
	if err != nil {
		log.Fatal(err)
	}

	w, err := smtpCon.Data()
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(msg)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = smtpCon.Quit()
	if err != nil {
		log.Fatal(err)
	}

	pop3Conn, err := net.Dial("tcp", "pop3.mailtrap.io:1100")
	if err != nil {
		log.Fatal(err)
	}

	c, err := pop3.NewClient(pop3Conn, pop3.UseTLS(conf), pop3.UseTimeout(time.Second*1))
	if err != nil {
		log.Fatal(err)
	}

	err = c.Auth(os.Getenv("MAILTRAP_USERNAME"), os.Getenv("MAILTRAP_PASSWORD"))
	if err != nil {
		log.Fatal(err)
	}

	count, size, err := c.Stat()
	if err != nil {
		log.Fatal(err)
	}

	data, err := c.ListAll()
	if err != nil {
		log.Fatal(err)
	}

	var messages []string
	for _, message := range data {
		fmt.Println(message.Seq)
		str, err := c.Retr(message.Seq)
		if err != nil {
			log.Fatal(err)
		}
		messages = append(messages, str)
		err = c.Dele(message.Seq)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = c.Quit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v, %v\n", count, size)
	for _, msg := range messages {
		fmt.Println(msg)
	}
}
