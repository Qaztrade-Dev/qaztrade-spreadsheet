package adapters

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"os"
	"testing"
)

//go:embed B.pdf
var documentPDF []byte

func TestStorageS3(t *testing.T) {
	var (
		ctx           = context.Background()
		accessKey     = os.Getenv("S3_ACCESS_KEY")
		secretKey     = os.Getenv("S3_SECRET_KEY")
		endpoint      = os.Getenv("S3_ENDPOINT")
		bucketName    = os.Getenv("S3_BUCKET")
		spreadsheetID = os.Getenv("TEMPLATE_SPREADSHEET_ID")
		sheetName     = "Затраты на доставку"

		fileName = "B.pdf"
		fileSize = int64(len(documentPDF))
		reader   = bytes.NewReader(documentPDF)
	)

	storage, err := NewStorageS3(ctx, accessKey, secretKey, bucketName, endpoint)
	if err != nil {
		t.Fatal("NewStorageS3 error:", err)
	}

	// 3. upload file, get url
	folderName := fmt.Sprintf("%s/%s", spreadsheetID, sheetName)

	value, err := storage.Upload(ctx, folderName, fileName, fileSize, reader)
	fmt.Println(value, err)
	if err != nil {
		t.Fatal("NewStorageS3 Upload error:", err)
	}
}
