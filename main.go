package main

import (
	"log"
	"net/smtp"
	"os"
	"strconv"
)

func main() {
	// Choose auth method and set it up
	host := "smtp.mailtrap.io"
	port := 25
	auth := smtp.PlainAuth("", os.Getenv("MAILTRAP_USERNAME"), os.Getenv("MAILTRAP_PASSWORD"), host)

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{"to@example.com"}
	msg := []byte("To: to@example.com\r\n" +
		"Subject: Why are you not using Mailtrap yet?\r\n" +
		"\r\n" +
		"Here’s the space for our great sales pitch\r\n")
	err := smtp.SendMail(host + ":" + strconv.Itoa(port), auth, "from@example.com", to, msg)
	if err != nil {
		log.Fatal(err)
	}
}
