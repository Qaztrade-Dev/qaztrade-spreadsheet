package service

import (
	"context"
	"fmt"
	"io"
	"time"
)

type SendNoticeRequest struct {
	ApplicationID string
	FileReader    io.Reader
	FileSize      int64
}

func (s *service) SendNotice(ctx context.Context, req *SendNoticeRequest) (string, error) {

	exists, errBucketExists := s.storage.BucketExists("qaztrade")
	fmt.Println(exists, errBucketExists)
	now := time.Now()
	folderName := "notices"
	fileName := fmt.Sprintf("%s-%s", req.ApplicationID, now.Format("2006-01-02-15-04"))
	value, err := s.storage.Upload(ctx, folderName, fileName, req.FileSize, req.FileReader)

	if err != nil {
		return "", err
	}
	return value, nil
}
