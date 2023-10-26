package doodocs

import (
	"context"
	"io"
	"time"
)

type CreateDocumentResponse struct {
	DocumentID string
	SignLink   string
}

type SigningService interface {
	GetSigningTime(ctx context.Context, documentID string) (time.Time, error)
	CreateDocument(ctx context.Context, documentName string, documentReader io.Reader) (*CreateDocumentResponse, error)
}
