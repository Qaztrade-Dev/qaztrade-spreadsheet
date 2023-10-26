package emailer

import (
	"context"
	"crypto/tls"
	_ "embed"
	"io"

	"gopkg.in/gomail.v2"
)

type EmailServiceGmail struct {
	email    string
	password string
	dialer   *gomail.Dialer
}

var _ EmailService = (*EmailServiceGmail)(nil)

func NewEmailerClient(email, password string) *EmailServiceGmail {

	dialer := gomail.NewDialer("webmail.p-s.kz", 587, email, password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &EmailServiceGmail{
		email:    email,
		password: password,
		dialer:   dialer,
	}
}

func (s *EmailServiceGmail) Send(ctx context.Context, mail *Email) error {
	message := gomail.NewMessage()
	message.SetHeader("From", s.email)
	message.SetHeader("To", mail.ToEmail)
	message.SetHeader("Subject", mail.Subject)
	message.SetBody("text/plain", mail.Body)

	if mail.AttachmentName != "" && mail.Attachment != nil {
		message.Attach(
			mail.AttachmentName,
			gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := io.Copy(w, mail.Attachment)
				return err
			}),
		)
	}

	if err := s.dialer.DialAndSend(message); err != nil {
		return err
	}

	return nil
}
