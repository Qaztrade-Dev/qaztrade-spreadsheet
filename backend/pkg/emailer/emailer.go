package emailer

import (
	"context"
	"io"
)

type Email struct {
	ToEmail        string
	Subject        string
	Body           string
	AttachmentName string
	Attachment     io.Reader
}

type EmailService interface {
	Send(ctx context.Context, mail *Email) error
}
