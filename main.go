package main

import (
	"log"
)

func main() {
	/*conf := NewConfig()
	errDB := CreateDatabase(conf)
	if err != nil {
		log.Fatal(err)
	}
	errDB = CreateTable(conf)
	if err != nil {
		log.Fatal(err)
	}*/
	s := NewServer()
	err := s.Start()
	if err != nil {
		log.Fatal(err)
	}
}
