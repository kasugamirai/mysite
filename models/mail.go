package models

import (
	"log"
	"net/smtp"
)

// mail config
var (
	smtpServer = "smtp.example.com"
	smtpPort   = "587"
	smtpUser   = "user@example.com"
	smtpPass   = "password"
)

type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func sendEmail(to string, subject string, body string) {
	from := smtpUser
	pass := smtpPass
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail(smtpServer+":"+smtpPort,
		smtp.PlainAuth("", from, pass, smtpServer),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent, visit http://foobar.com/baz")
}
