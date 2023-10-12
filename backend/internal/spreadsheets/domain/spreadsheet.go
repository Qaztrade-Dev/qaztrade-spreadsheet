package domain

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var ErrorSheetPresent = errors.New("sheet already present")

type SpreadsheetService interface {
	Create(ctx context.Context, user *User) (spreadsheetID string, err error)
	GetPublicLink(ctx context.Context, spreadsheetID string) (link string)
	AddSheet(ctx context.Context, spreadsheetID string, sheetName string) error
}

type Application struct {
	UserID        string
	SpreadsheetID string
	ApplicationNo int
	Link          string
	Status        string
	CreatedAt     time.Time
}

type ApplicationList struct {
	OverallCount uint64
	Applications []*Application
}

type ApplicationQuery struct {
	UserID string
	Limit  uint64
	Offset uint64
}

type ApplicationRepository interface {
	Create(ctx context.Context, userID string, input *Application) error
	GetMany(ctx context.Context, query *ApplicationQuery) (*ApplicationList, error)
}

type User struct {
	ID      string
	OrgName string
}

type UserRepository interface {
	Get(ctx context.Context, userID string) (*User, error)
}

func CreateSpreadsheetName(user *User) (string, error) {
	now := time.Now()

	location, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		return "", err
	}

	timeStr := now.In(location).Format(time.DateTime)
	return fmt.Sprintf("%s-%s-%s", user.ID, user.OrgName, timeStr), nil
}

type SpreadsheetClaims struct {
	SpreadsheetID string `json:"sid"`
}
