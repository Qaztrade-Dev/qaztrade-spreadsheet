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

	GetArchive(ctx context.Context, req *GetArchiveRequest) (*GetArchiveResponse, error)

	// EnqueueAssignments enqueues to check queue assignments that will be checked
	EnqueueAssignments(ctx context.Context) error

	// CheckAssignment checks the assignment
	CheckAssignment(ctx context.Context, assignmentID uint64) error
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
	// publisher domain.Publisher,
	spreadsheetRepo domain.SpreadsheetRepository,
) Service {
	return &service{
		assignmentRepo:  assignmentRepo,
		storage:         storage,
		spreadsheetRepo: spreadsheetRepo,
	}
}
