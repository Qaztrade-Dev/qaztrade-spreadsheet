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

func (s *service) SendNotice(ctx context.Context, req *SendNoticeRequest) error {
	application, err := s.applicationRepo.GetOne(ctx, &domain.GetManyInput{
		ApplicationID: req.ApplicationID,
	})
	if err != nil {
		return err
	}

	if application.Status != domain.StatusManagerReviewing {
		return domain.ErrorApplicationNotUnderReview
	}

	applicationAttr, err := s.spreadsheetSvc.GetApplication(ctx, application.SpreadsheetID)
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(nil)
	if _, err = io.Copy(buffer, req.FileReader); err != nil {
		return err
	}

	var (
		readerForStorage = bytes.NewReader(buffer.Bytes())
		readerForEmail   = bytes.NewReader(buffer.Bytes())

		now        = time.Now()
		folderName = "notices"
		fileName   = fmt.Sprintf("%s-%s", req.ApplicationID, now.Format("2006-01-02-15-04"))
	)

	_, err = s.storage.Upload(ctx, folderName, fileName, req.FileSize, readerForStorage)
	if err != nil {
		return err
	}

	if err := s.emailSvc.SendNotice(ctx, applicationAttr.ContEmail, "Уведомление по заявке от АО «QazTrade»", req.FileName, readerForEmail); err != nil {
		return err
	}

	if err := s.spreadsheetSvc.SwitchModeEdit(ctx, application.SpreadsheetID); err != nil {
		return err
	}

	if err := s.spreadsheetSvc.LockSheets(ctx, application.SpreadsheetID); err != nil {
		return err
	}

	if err := s.applicationRepo.EditStatus(ctx, req.ApplicationID, domain.StatusUserFixing); err != nil {
		return err
	}

	return nil
}
