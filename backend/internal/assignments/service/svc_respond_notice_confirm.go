package service

import (
	_ "embed"
	"time"

	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"golang.org/x/net/context"
)

func (s *service) RespondNoticeConfirm(ctx context.Context, documentID string) error {
	/*
		1. get message with doodocs_document_id
		2. get assignment by message.assignment_id
		3. check if assignment status is on fix
		4. check if assignment countdown duration is not over
		5. set assignment fields:
			- resolved_at = null
			- countdown_duration = null
			- resolution_status_id = manager_reviewing
		6. set message fields:
			- doodocs_signed_at = now
			- doodocs_is_signed = true
	*/

	message, err := s.msgRepo.GetOne(ctx, &domain.GetMessageInput{
		DoodocsDocumentID: documentID,
	})
	if err != nil {
		return err
	}

	assignment, err := s.assignmentRepo.GetOne(ctx, &domain.GetManyInput{
		AssignmentID: &message.AssignmentID,
	})
	if err != nil {
		return err
	}
	if err := checkAssignmentCanBeResponded(assignment); err != nil {
		return err
	}

	if err := s.assignmentRepo.SetResolution(ctx, &domain.SetResolutionInput{
		AssignmentID:      assignment.AssignmentID,
		ResolvedAt:        nil,
		CountdownDuration: nil,
		ResolutionStatus:  domain.ResolutionStatusOnReview,
	}); err != nil {
		return err
	}

	if err := s.msgRepo.UpdateMessage(ctx, &domain.UpdateMessageInput{
		MessageID:       message.MessageID,
		DoodocsSignedAt: time.Now().UTC(),
		DoodocsIsSigned: true,
	}); err != nil {
		return err
	}

	allReviewing, err := s.assignmentRepo.AllAssignmentsStatusEq(ctx, assignment.ApplicationID, domain.ResolutionStatusOnReview)
	if err != nil {
		return err
	}

	if allReviewing {
		if err := s.applicationRepo.EditStatus(ctx, assignment.ApplicationID, domain.ApplicationStatusOnReview); err != nil {
			return err
		}
	}

	return nil
}
