package email

import (
	"fmt"
	"net/smtp"
)

type Sender struct {
	Host string
	Port string
	User string
	Pass string
	From string
}

func NewSender(host, port, user, pass, from string) *Sender {
	return &Sender{Host: host, Port: port, User: user, Pass: pass, From: from}
}

func (s *Sender) Send(to, subject, body string) error {
	if s.User == "" || s.Pass == "" {
		return fmt.Errorf("SMTP not configured: set SMTP_USER and SMTP_PASS env vars")
	}
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.From, to, subject, body)
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)
	auth := smtp.PlainAuth("", s.User, s.Pass, s.Host)
	return smtp.SendMail(addr, auth, s.From, []string{to}, []byte(msg))
}
