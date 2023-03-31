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
}

var _ domain.EmailService = (*EmailServiceGmail)(nil)

func NewEmailServiceGmail(email, password string) *EmailServiceGmail {
	t := template.Must(template.ParseFiles("./templates/forgot.html"))
	return &EmailServiceGmail{
		email:    email,
		password: password,
		template: t,
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
	d := gomail.NewDialer("smtp.gmail.com", 587, s.email, s.password)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
