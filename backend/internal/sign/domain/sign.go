package domain

import (
	"bytes"
	"io"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
	"golang.org/x/net/context"
)

type PDFService interface {
	Create(application *Application, attachments []io.ReadSeeker) (*bytes.Buffer, error)
}

type SpreadsheetRepository interface {
	GetApplication(ctx context.Context, spreadsheetID string) (*domain.Application, error)
	GetAttachments(ctx context.Context, spreadsheetID string) ([]*bytes.Buffer, error)
}
