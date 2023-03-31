package domain

import (
	"context"
	"fmt"
	"time"
)

type SpreadsheetService interface {
	Create(ctx context.Context, user *User) (spreadsheetID string, err error)
	GetPublicLink(ctx context.Context, spreadsheetID string) (link string)
}

type Application struct {
	UserID        string
	SpreadsheetID string
	Link          string
}

type ApplicationRepository interface {
	Create(ctx context.Context, userID string, input *Application) error
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
