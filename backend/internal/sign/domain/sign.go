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

type CreateSigningDocumentResponse struct {
	DocumentID string
	SignLink   string
}

type SigningService interface {
	CreateSigningDocument(ctx context.Context, documentName string, documentReader io.Reader) (*CreateSigningDocumentResponse, error)
}

const (
	StatusUserFilling      = "user_filling"
	StatusManagerReviewing = "manager_reviewing"
	StatusUserFixing       = "user_fixing"
	StatusCompleted        = "completed"
	StatusRejected         = "rejected"
)

type SignApplication struct {
	SignLink string
	Status   string
}

type ApplicationRepository interface {
	AssignSigningInfo(ctx context.Context, spreadsheetID string, info *CreateSigningDocumentResponse) error
	EditStatus(ctx context.Context, spreadsheetID, statusName string) error
	GetApplication(ctx context.Context, spreadsheetID string) (*SignApplication, error)
}

type SpreadsheetClaims struct {
	SpreadsheetID string `json:"sid"`
}
