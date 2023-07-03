package domain

import (
	"context"
	"fmt"
	"time"
)

const (
	TypeDigital = "digital"
	TypeFinance = "finance"
	TypeLegal   = "legal"
)

type AssignmentInput struct {
	ApplicationID  string
	SheetTitle     string
	SheetID        uint64
	AssignmentType string
	ManagerID      string
	TotalRows      uint64
	TotalSum       float64
}

type AssignmentView struct {
	ID             int
	ApplicantName  string
	ApplicantBIN   string
	SheetTitle     string
	SheetID        uint64
	AssignmentType string
	Link           string
	AssigneeName   string
	TotalRows      int
	TotalSum       int
	RowsCompleted  int
	IsCompleted    bool
	CompletedAt    time.Time
}

type AssignmentsInfo struct {
	Total     uint64
	Completed uint64
}

type AssignmentsList struct {
	Total   int
	Objects []*AssignmentView
}

var (
	ErrAssignmentNotFound = fmt.Errorf("assignment not found")
)

type GetManyInput struct {
	UserID *string
	Limit  uint64
	Offset uint64
}

type GetInfoInput struct {
	UserID *string
}

var (
	ErrorEmptySheets   = fmt.Errorf("empty sheets")
	ErrorEmptyManagers = fmt.Errorf("empty managers")
)

type AssignmentsRepository interface {
	GetInfo(ctx context.Context, input *GetInfoInput) (*AssignmentsInfo, error)
	GetMany(ctx context.Context, input *GetManyInput) (*AssignmentsList, error)

	// LockApplications locks signed applications into a batch. Returns batch ID of newly created batch
	LockApplications(ctx context.Context) (int, error)

	// GetSheets returns sheets of a given sheet type
	GetSheets(ctx context.Context, batchID int, sheetTable string) ([]*Sheet, error)

	// GetManagerIDs returns ID of managers with specified role
	GetManagerIDs(ctx context.Context, role string) ([]string, error)

	// CreateAssignments creates given assignments
	CreateAssignments(ctx context.Context, inputs []*AssignmentInput) error

	// UpdateBatchStep update step of the given batch
	UpdateBatchStep(ctx context.Context, batchID, step int) error
}
