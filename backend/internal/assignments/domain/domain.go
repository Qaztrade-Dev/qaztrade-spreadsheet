package domain

import (
	"context"
	"fmt"
	"time"
)

type Sheet struct {
	SheetTitle string
	SheetID    int
	TotalRows  int
}

type Application struct {
	ApplicationID string
	SpreadsheetID string
	Sheets        []*Sheet
}

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
	Total     int
	Completed int
}

type AssignmentsList struct {
	Total   int
	Objects []*AssignmentView
}

var (
	ErrAssignmentNotFound = fmt.Errorf("assignment not found")
)

type AssignmentSearchInput struct {
	UserID string
	Limit  int
	Offset int
}

type InfoSearchInput struct {
	UserID string
}

type AssignmentRepository interface {
	GetInfo(ctx context.Context, input *InfoSearchInput) (*AssignmentsInfo, error)
	GetMany(ctx context.Context, input *AssignmentSearchInput) (*AssignmentsList, error)
}
