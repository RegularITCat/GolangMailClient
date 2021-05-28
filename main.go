package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/simia-tech/go-pop3"
	"io"
	"log"
	"net"
	"net/mail"
	//"net/smtp"
	"os"
	//"strconv"
	"strings"
	"time"
)

func CreateDatabase(DBName string) error {
	_, err := os.Create(DBName)
	if err != nil {
		return err
	}
	return nil
}

func CreateTable(DBName string) error {
	db, err := sql.Open("sqlite3", DBName)
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE `mails` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `from` TEXT, `to` TEXT, `subject` TEXT, `body` TEXT);")
	if err != nil {
		return err
	}
	err = db.Close()
	if err != nil {
		return err
	}
	return nil
}

/*func SendMessage() error {
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
	return nil
}
*/

func main() {

	conf := &tls.Config{ServerName: "pop3.mailtrap.io"}

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

	/*count, size, err := c.Stat()
	if err != nil {
		log.Fatal(err)
	}*/

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
	DBName := "./mails.db"
	/*err = CreateDatabase(DBName)
	if err != nil {
		log.Fatal(err)
	}
	err = CreateTable(DBName)
	if err != nil {
		log.Fatal(err)
	}*/
	db, err := sql.Open("sqlite3", DBName)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%v, %v\n", count, size)
	for _, msg := range messages {
		reader := strings.NewReader(msg)
		mailResponse, err := mail.ReadMessage(reader)
		if err != nil {
			log.Fatal(err)
		}
		header := mailResponse.Header
		body, err := io.ReadAll(mailResponse.Body)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(body)
		//`id` INTEGER PRIMARY KEY AUTOINCREMENT, `from` TEXT, `to` TEXT, `subject` TEXT, `body` TEXT
		query := "INSERT INTO \"mails\" (\"from\", \"to\", \"subject\", \"body\") values($1,$2,$3,$4);"
		result, err := db.Exec(query, header.Get("From"),
			header.Get("To"), header.Get("Subject"), body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result.LastInsertId())
	}
	err = db.Close()
	if err != nil {
		log.Fatal(err)
	}

}
