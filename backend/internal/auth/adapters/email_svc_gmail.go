package adapters

import (
	"bytes"
	"context"
	"crypto/tls"
	"text/template"

	"github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"gopkg.in/gomail.v2"
)

type EmailServiceGmail struct {
	email    string
	password string
	template *template.Template
	dialer   *gomail.Dialer
}

var _ domain.EmailService = (*EmailServiceGmail)(nil)

func NewEmailServiceGmail(email, password string) *EmailServiceGmail {
	templ := template.Must(template.ParseFiles("./templates/forgot.html"))

	dialer := gomail.NewDialer("webmail.p-s.kz", 587, email, password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &EmailServiceGmail{
		email:    email,
		password: password,
		template: templ,
		dialer:   dialer,
	}
}

func (s *EmailServiceGmail) Send(ctx context.Context, toEmail, mailName string, payload interface{}) error {
	buf := new(bytes.Buffer)
	err := s.template.Execute(buf, payload)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.email)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Восстановление пароля")
	m.SetBody("text/html", buf.String())

	if err := s.dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
