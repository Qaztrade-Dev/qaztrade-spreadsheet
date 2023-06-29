package domain

import (
	"context"
	"fmt"
	"time"
)

// type Sheet struct {
// 	SheetTitle string
// 	SheetID    int
// 	TotalRows  int
// }

// type Application struct {
// 	ApplicationID string
// 	SpreadsheetID string
// 	Sheets        []*Sheet
// }

const (
	TypeDigital = "digital"
	TypeFinance = "finance"
	TypeLegal   = "legal"
)

type AssignmentView struct {
	ID             int
	ApplicantName  string
	ApplicantBIN   string
	SheetTitle     string
	SheetID        uint64
	AssignmentType string
	Link           string
	AssigneeName   string
	RowsFrom       int
	RowsUntil      int
	RowsTotal      int
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

type AssignmentsRepository interface {
	GetInfo(ctx context.Context, input *GetInfoInput) (*AssignmentsInfo, error)
	GetMany(ctx context.Context, input *GetManyInput) (*AssignmentsList, error)
}
