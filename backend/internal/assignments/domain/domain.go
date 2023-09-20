package domain

import (
	"context"
	"fmt"
	"io"
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
	AssignmentID   uint64
	ID             uint64
	ApplicantName  string
	ApplicantBIN   string
	SpreadsheetID  string
	SheetTitle     string
	SheetID        uint64
	AssignmentType string
	Link           string
	SignLink       string
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

var (
	ErrorEmptySheets   = fmt.Errorf("empty sheets")
	ErrorEmptyManagers = fmt.Errorf("empty managers")
)

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
}

type RemoveFunction func() error

type Storage interface {
	GetArchive(ctx context.Context, folderName string) (io.ReadCloser, RemoveFunction, error)
}

type Publisher interface {
	Publish(ctx context.Context, assignmentID ...uint64) error
}

var (
	ErrorSheetNotFound = fmt.Errorf("sheet not found")
)

type SpreadsheetRepository interface {
	GetSheetData(ctx context.Context, spreadsheetID string, sheetTitle string) ([][]string, error)
}
