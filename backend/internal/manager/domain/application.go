package domain

import (
	"context"
	"time"
)

const (
	StatusUserFilling      = "user_filling"
	StatusManagerReviewing = "manager_reviewing"
	StatusUserFixing       = "user_fixing"
	StatusCompleted        = "completed"
	StatusRejected         = "rejected"
)

type Application struct {
	UserID        string
	SpreadsheetID string
	Link          string
	Status        string
	CreatedAt     time.Time
}

type ApplicationList struct {
	OverallCount uint64
	Applications []*Application
}

type ApplicationQuery struct {
	ApplicationID string
	Limit         uint64
	Offset        uint64
}

type ApplicationRepository interface {
	GetMany(ctx context.Context, query *ApplicationQuery) (*ApplicationList, error)
	GetOne(ctx context.Context, query *ApplicationQuery) (*Application, error)
	EditStatus(ctx context.Context, applicationID, statusName string) error
}

type SpreadsheetService interface {
	SwitchModeRead(ctx context.Context, spreadsheetID string) error
	SwitchModeEdit(ctx context.Context, spreadsheetID string) error
}
