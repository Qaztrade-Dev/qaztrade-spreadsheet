package service

import (
	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"golang.org/x/net/context"
)

type Service interface {
	// CreateBatch creates a batch of applications for assignment distribution.
	CreateBatch(ctx context.Context) error

	// GetUserAssignments returns assignments of a user.
	GetUserAssignments(ctx context.Context, input *GetUserAssignmentsRequest) (*GetUserAssignmentsResponse, error)

	// GetAssignments returns all assignments.
	GetAssignments(ctx context.Context, input *GetAssignmentsRequest) (*GetAssignmentsResponse, error)

	// ChangeAssignee changes assignee to the given assignee
	ChangeAssignee(ctx context.Context, input *ChangeAssigneeRequest) error
}

type service struct {
	assignmentRepo domain.AssignmentsRepository
}

func NewService(
	assignmentRepo domain.AssignmentsRepository,
) Service {
	return &service{
		assignmentRepo: assignmentRepo,
	}
}
