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

type GetUserAssignmentsInput struct {
	UserID string
	Limit  int
	Offset int
}

type GetUserAssignmentsOutput struct {
	AssignmentsList *domain.AssignmentsList
	AssignmentsInfo *domain.AssignmentsInfo
}

func (s *service) GetUserAssignments(ctx context.Context, input *GetUserAssignmentsInput) (*GetUserAssignmentsOutput, error) {
	assignmentsList, err := s.assignmentRepo.GetMany(ctx, &domain.AssignmentSearchInput{
		UserID: input.UserID,
		Limit:  input.Limit,
		Offset: input.Offset,
	})
	if err != nil {
		return nil, err
	}

	assignmentsInfo, err := s.assignmentRepo.GetInfo(ctx, &domain.InfoSearchInput{
		UserID: input.UserID,
	})
	if err != nil {
		return nil, err
	}

	output := &GetUserAssignmentsOutput{
		AssignmentsList: assignmentsList,
		AssignmentsInfo: assignmentsInfo,
	}

	return output, nil
}

type GetAssignmentsInput struct {
	Limit  int
	Offset int
}

type GetAssignmentsOutput struct {
	AssignmentsList *domain.AssignmentsList
	AssignmentsInfo *domain.AssignmentsInfo
}

func (s *service) GetAssignments(ctx context.Context, input *GetAssignmentsInput) (*GetAssignmentsOutput, error) {
	assignmentsList, err := s.assignmentRepo.GetMany(ctx, &domain.AssignmentSearchInput{
		Limit:  input.Limit,
		Offset: input.Offset,
	})
	if err != nil {
		return nil, err
	}

	assignmentsInfo, err := s.assignmentRepo.GetInfo(ctx, &domain.InfoSearchInput{})
	if err != nil {
		return nil, err
	}

	output := &GetAssignmentsOutput{
		AssignmentsList: assignmentsList,
		AssignmentsInfo: assignmentsInfo,
	}

	return output, nil
}

func NewService() Service {
	return &service{}
}
