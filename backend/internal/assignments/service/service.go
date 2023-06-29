package service

import (
	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"golang.org/x/net/context"
)

type Service interface {
	// CreateBatch creates a batch of applications for assignment distribution.
	CreateBatch(ctx context.Context) error

	// GetUserAssignments returns assignments of a user.
	GetUserAssignments(ctx context.Context, input *GetUserAssignmentsInput) (*GetUserAssignmentsOutput, error)

	// GetAssignments returns all assignments.
	GetAssignments(ctx context.Context, input *GetAssignmentsInput) (*GetAssignmentsOutput, error)
}

type service struct {
	assignmentRepo domain.AssignmentRepository
}

func (s *service) CreateBatch(ctx context.Context) error {
	// 1. get signed applications not in a batch
	// 2. create an empty batch (step=0)
	// 3. group signed applications into a single batch
	// 4. get a list of users (role=digital)
	// 5. distribute applications into assignments among these users
	// 6. create assignments
	// 7. start the batch (step=1)
	return nil
}

func NewService() Service {
	return &service{}
}
