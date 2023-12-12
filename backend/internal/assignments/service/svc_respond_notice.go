package service

import (
	_ "embed"
	"fmt"
	"io"

	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	applicationDomain "github.com/doodocs/qaztrade/backend/internal/manager/domain"

	"golang.org/x/net/context"
)

type RespondNoticeRequest struct {
	UserID         string
	ApplicationID  string
	AssignmentType string
	FileReader     io.Reader
	FileSize       int64
	FileName       string
}

type RespondNoticeResponse struct {
	SignLink string
}

func (s *service) RespondNotice(ctx context.Context, input *RespondNoticeRequest) (*RespondNoticeResponse, error) {
	/*
		0. get application by application_id
		1. check if user is author of the application
		2. get assignment by application_id and type
		3. check if assignment status is on fix
		4. check if assignment countdown duration is not over
		5. upload file to doodocs.kz and get document_id and link
		6. create new message with document_id
		7. return signlink
	*/

	application, err := s.applicationRepo.GetOne(ctx, &applicationDomain.GetManyInput{
		ApplicationID: input.ApplicationID,
	})
	if err != nil {
		return nil, err
	}

	if application.UserID != input.UserID {
		return nil, domain.ErrUnauthorized
	}

	assignment, err := s.assignmentRepo.GetOne(ctx, &domain.GetManyInput{
		ApplicationNo:  &application.No,
		AssignmentType: &input.AssignmentType,
	})
	if err != nil {
		return nil, err
	}

	if err := checkAssignmentCanBeResponded(assignment); err != nil {
		return nil, err
	}

	documentName := fmt.Sprintf("Сопровод:%s", input.FileName)

	doodocsResp, err := s.doodocs.CreateDocument(ctx, documentName, input.FileReader)
	if err != nil {
		return nil, err
	}

	if err := s.msgRepo.CreateMessage(ctx, &domain.CreateMessageInput{
		AssignmentID: assignment.AssignmentID,
		UserID:       input.UserID,
		Attrs: domain.MessageAttrs{
			"sign_link": doodocsResp.SignLink,
		},
		DoodocsDocumentID: doodocsResp.DocumentID,
	}); err != nil {
		return nil, err
	}

	return &RespondNoticeResponse{
		SignLink: doodocsResp.SignLink,
	}, nil
}

func checkAssignmentCanBeResponded(assignment *domain.AssignmentView) error {
	if assignment.ResolutionStatus != domain.ResolutionStatusOnFix {
		return domain.ErrAssignmentNotOnFix
	}

	// now := time.Now().UTC()

	// if assignment.ResolvedAt.UTC().Add(assignment.CountdownDuration).Before(now) {
	// 	return domain.ErrAssignmentCountdownDurationOver
	// }

	return nil
}
