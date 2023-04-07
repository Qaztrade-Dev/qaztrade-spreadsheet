package domain

import (
	"bytes"
	"io"

	"golang.org/x/net/context"
)

type PDFService interface {
	Create(application *Application, attachments []io.ReadSeeker) (*bytes.Buffer, error)
}

type SpreadsheetRepository interface {
	GetApplication(ctx context.Context, spreadsheetID string) (*Application, error)
	GetAttachments(ctx context.Context, spreadsheetID string) ([]io.ReadSeeker, error)
}
