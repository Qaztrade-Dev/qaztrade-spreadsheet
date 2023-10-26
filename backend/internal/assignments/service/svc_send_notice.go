package service

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"time"

	applicationDomain "github.com/doodocs/qaztrade/backend/internal/manager/domain"

	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"github.com/doodocs/qaztrade/backend/pkg/emailer"
	"golang.org/x/net/context"
)

//go:embed text_template.txt
var textTemplate []byte

type SendNoticeRequest struct {
	UserID       string
	AssignmentID uint64
	FileReader   io.Reader
	FileSize     int64
	FileName     string
}

func (s *service) SendNotice(ctx context.Context, input *SendNoticeRequest) error {
	assignment, err := s.assignmentRepo.GetOne(ctx, &domain.GetManyInput{
		AssignmentID: &input.AssignmentID,
	})
	if err != nil {
		return err
	}

	application, err := s.applicationRepo.GetOne(ctx, &applicationDomain.GetManyInput{
		ApplicationID: assignment.ApplicationID,
	})
	if err != nil {
		return err
	}

	applicationAttr, err := s.spreadsheetRepo.GetApplicationAttrs(ctx, application.SpreadsheetID)
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(nil)
	if _, err = io.Copy(buffer, input.FileReader); err != nil {
		return err
	}

	var (
		readerForStorage = bytes.NewReader(buffer.Bytes())
		readerForEmail   = bytes.NewReader(buffer.Bytes())

		now        = time.Now().UTC()
		folderName = "notices"
		fileKey    = fmt.Sprintf(
			"%s/%s:%s:%s",
			folderName,
			assignment.ApplicationID,
			assignment.AssignmentType,
			now.Format(time.RFC3339),
		)
	)

	fileURL, err := s.storage.Upload(ctx, fileKey, input.FileSize, readerForStorage)
	if err != nil {
		return err
	}

	if err := s.assignmentRepo.SetResolution(ctx, &domain.SetResolutionInput{
		AssignmentID:      input.AssignmentID,
		CountdownDuration: &domain.DefaultCountdownDuration,
		ResolvedAt:        &now,
		ResolutionStatus:  domain.ResolutionStatusOnFix,
	}); err != nil {
		return err
	}

	if err := s.msgRepo.CreateMessage(ctx, &domain.CreateMessageInput{
		AssignmentID: input.AssignmentID,
		UserID:       input.UserID,
		Attrs: domain.MessageAttrs{
			"file_url": fileURL,
		},
	}); err != nil {
		return err
	}

	if err := s.spreadsheetRepo.SwitchModeEdit(ctx, application.SpreadsheetID); err != nil {
		return err
	}

	if err := s.spreadsheetRepo.LockSheets(ctx, application.SpreadsheetID); err != nil {
		return err
	}

	if err := s.applicationRepo.EditStatus(ctx, application.ID, domain.ApplicationStatusOnFix); err != nil {
		return err
	}

	if err := s.emailer.Send(ctx, &emailer.Email{
		ToEmail:        applicationAttr.ContEmail,
		Subject:        "Уведомление по заявке от АО «QazTrade»",
		Body:           string(textTemplate),
		AttachmentName: input.FileName,
		Attachment:     readerForEmail,
	}); err != nil {
		return err
	}

	return nil
}
