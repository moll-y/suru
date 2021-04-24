package main

import (
	"crypto/tls"
	"log"

	gomail "gopkg.in/mail.v2"
)

type MessageBroker struct {
	username string
	password string
	to       string
	host     string
	port     int
	dialer   *gomail.Dialer
	message  *gomail.Message
	messages chan string
}

func NewMessageBrokerService(host string, port int, username, password, to string) *MessageBroker {
	dialer := gomail.NewDialer(host, port, username, password)
	dialer.TLSConfig = &tls.Config{ServerName: host, InsecureSkipVerify: false}
	messages := make(chan string)
	message := gomail.NewMessage()
	message.SetHeader("From", username)
	return &MessageBroker{
		username: username,
		password: password,
		host:     host,
		port:     port,
		to:       to,
		dialer:   dialer,
		messages: messages,
		message:  message,
	}
}

func (m *MessageBroker) SendMessage(msg string) {
	m.messages <- msg
}

func (m *MessageBroker) Dispatcher() {
	for {
		select {
		case msg := <-m.messages:
			m.dispatch(msg)
		}
	}
}

func (m *MessageBroker) dispatch(msg string) {
	log.Println("sending an email.")
	m.message.SetHeader("To", m.to)
	m.message.SetHeader("Subject", msg)
	m.message.SetBody("text/plain", msg)
	if err := m.dialer.DialAndSend(m.message); err != nil {
		log.Printf(err.Error())
	}
}
