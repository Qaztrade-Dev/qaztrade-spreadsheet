package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
)

type SendNoticeRequest struct {
	ApplicationID string
	FileReader    io.Reader
	FileSize      int64
	FileName      string
}

func (s *service) SendNotice(ctx context.Context, req *SendNoticeRequest) (string, error) {

	application, err := s.applicationRepo.GetOne(ctx, &domain.ApplicationQuery{
		ApplicationID: req.ApplicationID,
	})

	if err != nil {
		return "", err
	}

	if application.Status != domain.StatusManagerReviewing {
		return "", domain.ErrorApplicationNotUnderReview
	}

	applicationAttr, err := s.spreadsheetSvc.GetApplication(ctx, application.SpreadsheetID)
	if err != nil {
		return "", err
	}
	buffer := bytes.Buffer{}
	_, err = io.Copy(&buffer, req.FileReader)
	if err != nil {
		return "", err
	}
	readSeeker := bytes.NewReader(buffer.Bytes())
	readSeeker2 := bytes.NewReader(buffer.Bytes())

	var (
		now        = time.Now()
		folderName = "notices"
		fileName   = fmt.Sprintf("%s-%s", req.ApplicationID, now.Format("2006-01-02-15-04"))
	)
	value, err := s.storage.Upload(ctx, folderName, fileName, req.FileSize, readSeeker)
	if err != nil {
		return "", err
	}

	if err := s.emailSvc.SendNotice(ctx, applicationAttr.ContEmail, "Уведомления по замечанием в заявление", req.FileName, readSeeker2); err != nil {
		return "", err
	}

	if err := s.applicationRepo.EditStatus(ctx, req.ApplicationID, domain.StatusUserFixing); err != nil {
		return "", err
	}

	return value, nil
}
