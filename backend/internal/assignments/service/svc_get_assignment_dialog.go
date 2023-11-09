package service

import (
	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"golang.org/x/net/context"
)

type GetAssignmentDialogResponse struct {
	Messages []*domain.Message
}

func (s *service) GetAssignmentDialog(ctx context.Context, assignmentID uint64) (*GetAssignmentDialogResponse, error) {
	messages, err := s.msgRepo.GetMany(ctx, &domain.GetMessageInput{
		AssignmentID: assignmentID,
	})
	if err != nil {
		return nil, err
	}

	return &GetAssignmentDialogResponse{
		Messages: messages,
	}, nil
}
