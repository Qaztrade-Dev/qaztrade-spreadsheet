package domain

import (
	"context"
	"fmt"
	"net/http"
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
	ID             string
	UserID         string
	No             int
	SpreadsheetID  string
	Link           string
	Status         string
	SignDocumentID string
	Attrs          interface{}
	SignedAt       time.Time
	CreatedAt      time.Time
}

type ApplicationList struct {
	OverallCount uint64
	Applications []*Application
}

type ApplicationQuery struct {
	Limit  uint64
	Offset uint64

	ApplicationID    string
	BIN              string
	CompensationType string
	SignedAtFrom     time.Time
	SignedAtUntil    time.Time
}

type Revision struct {
	ApplicationID  string
	SpreadsheetID  string
	No             int
	Link           string
	BIN            string
	Manufactor     string
	To             string
	ApplicantEmail string
	ManagerName    string
	ManagerEmail   string
	Remarks        string
}

type ApplicationRepository interface {
	GetMany(ctx context.Context, query *ApplicationQuery) (*ApplicationList, error)
	GetOne(ctx context.Context, query *ApplicationQuery) (*Application, error)
	EditStatus(ctx context.Context, applicationID, statusName string) error
}

type SpreadsheetService interface {
	SwitchModeRead(ctx context.Context, spreadsheetID string) error
	SwitchModeEdit(ctx context.Context, spreadsheetID string) error
	Comments(ctx context.Context, application *Application) (*Revision, error)
}

var (
	ErrorApplicationNotSigned = fmt.Errorf("Заявление еще не подписано!")
)

type SigningService interface {
	GetDDCard(ctx context.Context, documentID string) (*http.Response, error)
}
