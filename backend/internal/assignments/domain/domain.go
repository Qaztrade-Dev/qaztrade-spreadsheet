package domain

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrUnauthorized                    = errors.New("Пользователь не имеет доступа к данной заявке")
	ErrAssignmentNotOnFix              = errors.New("Задача не на исправлении")
	ErrAssignmentCountdownDurationOver = errors.New("Время на исправление задачи истекло")
)

const (
	TypeDigital = "digital"
	TypeFinance = "finance"
	TypeLegal   = "legal"

	ResolutionStatusOnReview  = "manager_reviewing"
	ResolutionStatusOnFix     = "user_fixing"
	ResolutionStatusCompleted = "completed"
	ResolutionStatusRejected  = "rejected"

	ApplicationStatusOnReview = "manager_reviewing"
	ApplicationStatusOnFix    = "user_fixing"
)

var (
	DefaultCountdownDuration = 7 * 24 * time.Hour
)

type AssignmentInput struct {
	AssignmentID   uint64
	ApplicationID  string
	SheetTitle     string
	SheetID        uint64
	AssignmentType string
	ManagerID      string
	TotalRows      uint64
	TotalSum       float64
}

type AssignmentView struct {
	ApplicationID     string
	AssignmentID      uint64
	ID                uint64
	ApplicantName     string
	ApplicantBIN      string
	SpreadsheetID     string
	SheetTitle        string
	SheetID           uint64
	AssignmentType    string
	Link              string
	SignLink          string
	AssigneeName      string
	AssigneeID        string
	TotalRows         int
	TotalSum          int
	RowsCompleted     int
	IsCompleted       bool
	CompletedAt       time.Time
	ResolutionStatus  string
	ResolvedAt        time.Time
	CountdownDuration time.Duration
}

type AssignmentsInfo struct {
	Total     uint64
	Completed uint64
}

type AssignmentsList struct {
	Total   int
	Objects []*AssignmentView
}

type GetManyInput struct {
	AssigneeID     *string
	AssignmentID   *uint64
	IsCompleted    *bool
	CompanyName    *string // название компании
	ApplicationNo  *int    // номер заявки
	AssignmentType *string // тип задачи
	Limit          uint64
	Offset         uint64
}

type GetInfoInput struct {
	UserID *string
}

type ChangeAssigneeInput struct {
	UserID       string
	AssignmentID uint64
}

type SetResolutionInput struct {
	AssignmentID      uint64
	CountdownDuration *time.Duration
	ResolvedAt        *time.Time
	ResolutionStatus  string
}

type AssignmentsRepository interface {
	GetInfo(ctx context.Context, input *GetInfoInput) (*AssignmentsInfo, error)
	GetMany(ctx context.Context, input *GetManyInput) (*AssignmentsList, error)
	GetOne(ctx context.Context, input *GetManyInput) (*AssignmentView, error)

	// LockApplications locks signed applications into a batch. Returns batch ID of newly created batch
	LockApplications(ctx context.Context) (int, error)

	// GetSheets returns sheets of a given sheet type
	GetSheets(ctx context.Context, batchID int, sheetTable string) ([]*Sheet, error)

	// GetManagerIDs returns ID of managers with specified role
	GetManagerIDs(ctx context.Context, role string) ([]string, error)

	// CreateAssignments creates given assignments
	CreateAssignments(ctx context.Context, inputs []*AssignmentInput) error

	// ChangeAssignee changes assignment assinee
	ChangeAssignee(ctx context.Context, input *ChangeAssigneeInput) error

	// InsertAssignmentResult inserts assignment result and updates along related tables
	InsertAssignmentResult(ctx context.Context, assignmentID uint64, total uint64) error

	UpdateAssignees(ctx context.Context, inputs []*AssignmentInput) error

	SetResolution(ctx context.Context, input *SetResolutionInput) error

	AllAssignmentsStatusEq(ctx context.Context, applicationID, statusName string) (bool, error)
}

var (
	ErrAssignmentNotFound = fmt.Errorf("assignment not found")
	ErrorEmptySheets      = fmt.Errorf("empty sheets")
	ErrorEmptyManagers    = fmt.Errorf("empty managers")
)

type Publisher interface {
	Publish(ctx context.Context, assignmentID ...uint64) error
}

type MessageAttrs map[string]interface{}

type CreateMessageInput struct {
	AssignmentID      uint64
	UserID            string
	Attrs             MessageAttrs
	DoodocsDocumentID string
}

type GetMessageInput struct {
	DoodocsDocumentID string
}

type Message struct {
	MessageID    string
	AssignmentID uint64
}

type UpdateMessageInput struct {
	MessageID       string
	DoodocsSignedAt time.Time
	DoodocsIsSigned bool
}

type MessagesRepository interface {
	CreateMessage(ctx context.Context, input *CreateMessageInput) error
	GetOne(ctx context.Context, input *GetMessageInput) (*Message, error)
	UpdateMessage(ctx context.Context, input *UpdateMessageInput) error
}

var (
	ErrMessageNotFound = fmt.Errorf("message not found")
)
