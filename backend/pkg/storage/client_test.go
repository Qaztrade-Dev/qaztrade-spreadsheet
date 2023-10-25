package storage

import (
	"context"
	_ "embed"
	"io"
	"os"
	"testing"
)

func TestStorageS3(t *testing.T) {
	var (
		ctx           = context.Background()
		accessKey     = os.Getenv("S3_ACCESS_KEY")
		secretKey     = os.Getenv("S3_SECRET_KEY")
		endpoint      = os.Getenv("S3_ENDPOINT")
		bucketName    = os.Getenv("S3_BUCKET")
		spreadsheetID = "1_77bLWWtHUHKkiwYoIqJGoOW7rTpGpCUCMdhj1eCERI"
	)

	storage, err := NewStorageS3(ctx, accessKey, secretKey, bucketName, endpoint)
	if err != nil {
		t.Fatal("NewStorageS3 error:", err)
	}

	readCloser, remove, err := storage.GetArchive(ctx, spreadsheetID)
	if err != nil {
		t.Fatal("NewStorageS3 Upload error:", err)
	}
	defer remove()

	b, err := io.ReadAll(readCloser)
	if err != nil {
		t.Fatal("io.ReadAll error:", err)
	}

	os.WriteFile("archive.zip", b, 0644)
}
