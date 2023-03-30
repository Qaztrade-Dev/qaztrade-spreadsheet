package adapters

import (
	"context"
	"testing"
)

func TestStorageS3(t *testing.T) {
	var (
		ctx       = context.Background()
		accessKey = "KZTX0SSBNETFJ2OM84VE"
		secretKey = "wqb0TRw0rIZyw8wQGejISZ6VOldAnpAMGFgNEC1U"
		endpoint  = "https://object.pscloud.io"
	)

	_, err := NewStorageS3(ctx, accessKey, secretKey, endpoint)
	if err != nil {
		t.Fatal("NewStorageS3 error:", err)
	}
}
