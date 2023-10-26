package jsonmanager

import (
	"time"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
)

type Application struct {
	ID            string      `json:"id"`
	No            int         `json:"no"`
	SpreadsheetID string      `json:"spreadsheet_id"`
	Link          string      `json:"link"`
	Status        string      `json:"status,omitempty"`
	CreatedAt     time.Time   `json:"created_at,omitempty"`
	SignedAt      time.Time   `json:"signed_at,omitempty"`
	Attrs         interface{} `json:"attrs"`
	TotalRows     int         `json:"total_rows"`
	TotalSum      int         `json:"total_sum"`
	AttrsDigital  interface{} `json:"digital,omitempty"`
	AttrsFinance  interface{} `json:"finance"`
	AttrsLegal    interface{} `json:"legal"`
}

func EncodeApplication(input *domain.Application) *Application {
	if input == nil {
		return nil
	}

	return &Application{
		ID:            input.ID,
		No:            input.No,
		SpreadsheetID: input.SpreadsheetID,
		Link:          input.Link,
		Status:        input.Status,
		CreatedAt:     input.CreatedAt,
		SignedAt:      input.SignedAt,
		Attrs:         input.Attrs,
		TotalRows:     input.TotalRows,
		TotalSum:      input.TotalSum,
		AttrsDigital:  input.AttrsDigital,
		AttrsFinance:  input.AttrsFinance,
		AttrsLegal:    input.AttrsLegal,
	}
}

type ApplicationList struct {
	OverallCount uint64         `json:"overall_count"`
	Applications []*Application `json:"applications"`
}

func EncodeApplicationList(input *domain.ApplicationList) *ApplicationList {
	if input == nil {
		return nil
	}

	return &ApplicationList{
		OverallCount: input.OverallCount,
		Applications: EncodeSlice(input.Applications, EncodeApplication),
	}
}

type Manager struct {
	UserID   string   `json:"user_id"`
	Email    string   `json:"email"`
	Fullname string   `json:"full_name"`
	Roles    []string `json:"roles"`
}

func EncodeManager(input *domain.Manager) *Manager {
	if input == nil {
		return nil
	}

	return &Manager{
		UserID:   input.UserID,
		Email:    input.Email,
		Fullname: input.Fullname,
		Roles:    input.Roles,
	}
}

// encoder func for encoding from D (domain) to J (json)
type encoder[D any, J any] func(*D) *J

// EncodeSlice encodes slice of domain objects to a slice of json objects.
func EncodeSlice[D any, J any](slice []*D, encode encoder[D, J]) []*J {
	result := make([]*J, len(slice))
	for i, v := range slice {
		result[i] = encode(v)
	}
	return result
}
