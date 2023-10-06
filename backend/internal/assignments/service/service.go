package service

import (
	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"golang.org/x/net/context"
)

type Service interface {
	// CreateBatch creates a batch of applications for assignment distribution.
	CreateBatch(ctx context.Context) error

	// GetUserAssignments returns assignments of a user.
	GetUserAssignments(ctx context.Context, input *domain.GetManyInput) (*GetUserAssignmentsResponse, error)

	// GetAssignments returns all assignments.
	GetAssignments(ctx context.Context, input *domain.GetManyInput) (*GetAssignmentsResponse, error)

	// ChangeAssignee changes assignee to the given assignee
	ChangeAssignee(ctx context.Context, input *ChangeAssigneeRequest) error

	GetArchive(ctx context.Context, req *GetArchiveRequest) (*GetArchiveResponse, error)

	// EnqueueAssignments enqueues to check queue assignments that will be checked
	EnqueueAssignments(ctx context.Context) error

	// CheckAssignment checks the assignment
	CheckAssignment(ctx context.Context, assignmentID uint64) error

	// RedistributeAssignments redistributes assignments
	RedistributeAssignments(ctx context.Context, assignmentType string) error
}

type service struct {
	assignmentRepo  domain.AssignmentsRepository
	storage         domain.Storage
	publisher       domain.Publisher
	spreadsheetRepo domain.SpreadsheetRepository
}

func NewService(
	assignmentRepo domain.AssignmentsRepository,
	storage domain.Storage,
	spreadsheetRepo domain.SpreadsheetRepository,
	publisher domain.Publisher,
) Service {
	return &service{
		assignmentRepo:  assignmentRepo,
		storage:         storage,
		spreadsheetRepo: spreadsheetRepo,
		publisher:       publisher,
	}
}
