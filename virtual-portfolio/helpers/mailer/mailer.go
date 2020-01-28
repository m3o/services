package mailer

import (
	"fmt"
	"net/smtp"
)

// smtpServer data to smtp server.
type smtpServer struct {
	host, port string
}

func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

// Service is an object which can send emailss
type Service interface {
	Send(Email) error
}

// Email is the message being sent
type Email struct {
	ToAddress string
	ToName    string
	Subject   string
	Body      string
}

// Handler is an implementation of Service
type Handler struct {
	email, password string
	server          smtpServer
	auth            smtp.Auth
}

// Send submits the Email to the mail serve
func (h Handler) Send(email Email) error {
	subject := fmt.Sprintf("Subject: %v\n", email.Subject)
	to := fmt.Sprintf("To: %v <%v>\n", email.ToName, email.ToAddress)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(subject + to + mime + "" + email.Body)

	return smtp.SendMail(h.server.Address(), h.auth, h.email, []string{email.ToAddress}, msg)
}

// New returns an instance of handler, given the credentials
func New(email, password string) Handler {
	server := smtpServer{host: "smtp.gmail.com", port: "587"}
	auth := smtp.PlainAuth("", email, password, server.host)
	return Handler{email, password, server, auth}
}
