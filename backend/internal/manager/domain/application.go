package domain

import (
	"context"
	"errors"
	"io"
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
	ApplicationID string
	Limit         uint64
	Offset        uint64

	BIN              string
	CompensationType string
	SignedAtFrom     time.Time
	SignedAtUntil    time.Time
}

type ApplicationRepository interface {
	GetMany(ctx context.Context, query *ApplicationQuery) (*ApplicationList, error)
	GetOne(ctx context.Context, query *ApplicationQuery) (*Application, error)
	EditStatus(ctx context.Context, applicationID, statusName string) error
}

type SpreadsheetService interface {
	SwitchModeRead(ctx context.Context, spreadsheetID string) error
	SwitchModeEdit(ctx context.Context, spreadsheetID string) error
	BlockImportantRanges(ctx context.Context, spreadsheetID string) error
	UnlockImportantRanges(ctx context.Context, spreadsheetID string) error
}

type RemoveFunction func() error

type Storage interface {
	DownloadArchive(ctx context.Context, folderName string) (io.ReadCloser, RemoveFunction, error)
}

var (
	ErrorApplicationNotSigned = errors.New("Заявление еще не подписано!")
)

type SigningService interface {
	GetDDCardResponse(ctx context.Context, documentID string) (*http.Response, error)
}
