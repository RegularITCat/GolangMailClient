package main

import (
	"log"
)

func main() {
	s := NewServer()
	err := s.Start()
	if err != nil {
		log.Fatal(err)
	}
}
