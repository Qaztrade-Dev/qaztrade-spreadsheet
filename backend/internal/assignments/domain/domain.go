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

type Assignment struct {
	ID            int
	UserID        int
	ApplicationID string
	SpreadsheetID string
	Type          string
	SheetTitle    string
	SheetID       int
	RowsFrom      int
	RowsUntil     int
	RowsTotal     int
	RowsCompleted int
	IsCompleted   bool
	CompletedAt   time.Time
}

type AssignmentsInfo struct {
	Total     int
	Completed int
}

type AssignmentsList struct {
	Total   int
	Objects []*Assignment
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
