package main

//TODO Make a method for mail to send himself lol!
type Mail struct {
	Id       int    `json:"id"`
	From     string `json:"from"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	FullText string `json:"full_text"`
}
