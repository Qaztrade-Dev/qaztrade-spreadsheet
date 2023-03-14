package storage

import (
	"context"
	"io"

	"github.com/pkg/errors"
)

type Service interface {
	UploadFile(ctx context.Context, req *UploadFileRequest) (*UploadFileResponse, error)
}

type UploadFileRequest struct {
	FileName      string
	DirectoryName string
	Reader        io.Reader
	FileSize      int64
}

type UploadFileResponse struct {
	FileName string
	Err      error
}

var (
	ErrorInvalidFileExtension = errors.New("invalid file extension")
)

type Storage interface {
	Upload(ctx context.Context, file io.Reader, size int64, name string) error
}
