package storage

import (
	"context"
	"io"
)

type RemoveFunction func() error

type Storage interface {
	Upload(ctx context.Context, filekey string, fileSize int64, fileReader io.Reader) (string, error)
	Remove(ctx context.Context, filePath string) error
	GetArchive(ctx context.Context, folderName string) (io.ReadCloser, RemoveFunction, error)
}
