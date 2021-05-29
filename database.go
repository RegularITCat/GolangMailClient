package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

func CreateDatabase(c *Config) error {
	_, err := os.Create(c.DBPath)
	if err != nil {
		return err
	}
	return nil
}

func CreateTable(c *Config) error {
	db, err := sql.Open(c.DBDriver, c.DBPath)
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE `mails` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `from` TEXT, `to` TEXT, `subject` TEXT, `fullText` TEXT);")
	if err != nil {
		return err
	}
	err = db.Close()
	if err != nil {
		return err
	}
	return nil
}

type MailMap struct {
	DBPath   string
	DBDriver string
}

func NewMailMap(c *Config) *MailMap {
	return &MailMap{
		DBPath:   c.DBPath,
		DBDriver: c.DBDriver,
	}
}

func (mm MailMap) Insert(m *Mail) error {
	db, err := sql.Open(mm.DBDriver, mm.DBPath)
	if err != nil {
		return err
	}
	query := "INSERT INTO \"mails\" (\"from\", \"to\", \"subject\", \"fullText\") values($1,$2,$3,$4);"
	result, err := db.Exec(query, m.From, m.To, m.Subject, m.FullText)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	err = db.Close()
	if err != nil {
		return err
	}
	m.Id = int(id)
	return nil
}

func (mm MailMap) Update(m Mail) error {
	db, err := sql.Open(mm.DBDriver, mm.DBPath)
	if err != nil {
		return err
	}
	query := "UPDATE \"mails\" SET \"from\" = $1, \"to\" = $2, \"subject\" = $3, \"fullText\" = $4 WHERE \"id\" = $5"
	_, err = db.Exec(query, m.From, m.To, m.Subject, m.FullText, m.Id)
	if err != nil {
		return err
	}
	err = db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (mm MailMap) Select(id uint) (Mail, error) {
	db, err := sql.Open(mm.DBDriver, mm.DBPath)
	if err != nil {
		return Mail{}, err
	}
	query := "select * from mails where id = $1"
	row, err := db.Query(query, id)
	if err != nil {
		return Mail{}, err
	}
	err = db.Close()
	if err != nil {
		return Mail{}, err
	}

	var mail Mail
	err = row.Scan(&mail.Id, &mail.From, &mail.To, &mail.Subject, &mail.FullText)
	if err != nil {
		return Mail{}, err
	}
	return mail, nil
}

func (mm *MailMap) SelectAll() (map[int]Mail, error) {
	mails := make(map[int]Mail)
	db, err := sql.Open(mm.DBDriver, mm.DBPath)
	if err != nil {
		return map[int]Mail{}, err
	}
	query := "select * from mails"
	rows, err := db.Query(query)
	if err != nil {
		return map[int]Mail{}, err
	}
	err = db.Close()
	if err != nil {
		return map[int]Mail{}, err
	}
	for rows.Next() {
		m := Mail{}
		err := rows.Scan(&m.Id, &m.From, &m.To, &m.Subject, &m.FullText)
		if err != nil {
			return map[int]Mail{}, err
		}
		mails[m.Id] = m
	}
	return mails, nil
}
