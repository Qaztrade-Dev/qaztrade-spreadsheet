package adapters

import (
	"context"
	"testing"
)

func TestStorageS3(t *testing.T) {
	var (
		ctx        = context.Background()
		accessKey  = "x"
		secretKey  = "x"
		endpoint   = "x"
		bucketName = "x"
	)

	_, err := NewStorageS3(ctx, accessKey, secretKey, bucketName, endpoint)
	if err != nil {
		t.Fatal("NewStorageS3 error:", err)
	}
}
