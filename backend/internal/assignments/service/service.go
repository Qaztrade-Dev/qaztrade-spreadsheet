package service

import (
	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	applicationDomain "github.com/doodocs/qaztrade/backend/internal/manager/domain"
	"github.com/doodocs/qaztrade/backend/pkg/doodocs"
	"github.com/doodocs/qaztrade/backend/pkg/emailer"
	"github.com/doodocs/qaztrade/backend/pkg/publisher"
	"github.com/doodocs/qaztrade/backend/pkg/spreadsheets"
	"github.com/doodocs/qaztrade/backend/pkg/storage"
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

	SendNotice(ctx context.Context, req *SendNoticeRequest) error

	RespondNotice(ctx context.Context, input *RespondNoticeRequest) (*RespondNoticeResponse, error)

	RespondNoticeConfirm(ctx context.Context, documentID string) error
}

type service struct {
	assignmentRepo  domain.AssignmentsRepository
	storage         storage.Storage
	spreadsheetRepo spreadsheets.SpreadsheetService
	publisher       publisher.Publisher
	emailer         emailer.EmailService
	msgRepo         domain.MessagesRepository
	applicationRepo applicationDomain.ApplicationRepository
	doodocs         doodocs.SigningService
}

func NewService(
	assignmentRepo domain.AssignmentsRepository,
	storage storage.Storage,
	spreadsheetRepo spreadsheets.SpreadsheetService,
	publisher publisher.Publisher,
	emailer emailer.EmailService,
	msgRepo domain.MessagesRepository,
	applicationRepo applicationDomain.ApplicationRepository,
	doodocs doodocs.SigningService,
) Service {
	return &service{
		assignmentRepo:  assignmentRepo,
		storage:         storage,
		spreadsheetRepo: spreadsheetRepo,
		publisher:       publisher,
		emailer:         emailer,
		msgRepo:         msgRepo,
		applicationRepo: applicationRepo,
		doodocs:         doodocs,
	}
}
