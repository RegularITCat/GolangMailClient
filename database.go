package main

import (
	"database/sql"
	"os"
)

func CreateDatabase(c Config) error {
	_, err := os.Create(c.DBPath)
	if err != nil {
		return err
	}
	return nil
}

func CreateTable(c Config) error {
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

func NewMailMap(c Config) *MailMap {
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
	result, err := db.Exec(query, m.from, m.to, m.subject, m.fullText)
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
	m.id = int(id)
	return nil
}

func (mm MailMap) Update(m Mail) error {
	db, err := sql.Open(mm.DBDriver, mm.DBPath)
	if err != nil {
		return err
	}
	query := "UPDATE \"mails\" SET \"from\" = $1, \"to\" = $2, \"subject\" = $3, \"fullText\" = $4 WHERE \"id\" = $5"
	_, err = db.Exec(query, m.from, m.to, m.subject, m.fullText, m.id)
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
	err = row.Scan(&mail.id, &mail.from, &mail.to, &mail.subject, &mail.fullText)
	if err != nil {
		return Mail{}, err
	}
	return mail, nil
}

func (mm *MailMap) SelectAll() ([]Mail, error) {
	var mails []Mail
	db, err := sql.Open(mm.DBDriver, mm.DBPath)
	if err != nil {
		return []Mail{}, err
	}
	query := "select * from mails"
	rows, err := db.Query(query)
	if err != nil {
		return []Mail{}, err
	}
	err = db.Close()
	if err != nil {
		return []Mail{}, err
	}
	for rows.Next() {
		m := Mail{}
		err := rows.Scan(&m.id, &m.from, &m.to, &m.subject, &m.fullText)
		if err != nil {
			return []Mail{}, err
		}
		mails = append(mails, m)
	}
	return mails, nil
}
