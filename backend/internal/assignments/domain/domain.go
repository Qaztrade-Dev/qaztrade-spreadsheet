package domain

import (
	"context"
	"fmt"
	"time"
)

type Assignment struct {
	UserID        int
	ApplicationID string
	SpreadsheetID string
	Type          string
	SheetTitle    string
	SheetID       int
	RowsFrom      int
	RowsUntil     int
	RowsTotal     int
	IsCompleted   bool
	CompletedAt   time.Time
}

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

type AssignmentSearchInput struct {
	UserID string
	Limit  int
	Offset int
}

var (
	ErrAssignmentNotFound = fmt.Errorf("assignment not found")
)

type AssignmentRepository interface {
	GetMany(ctx context.Context, input *AssignmentSearchInput) ([]*Assignment, error)
}
