package jsonspreadsheets

import (
	"time"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
)

type Application struct {
	ID            string    `json:"application_id"`
	SpreadsheetID string    `json:"spreadsheet_id"`
	ApplicationNo int       `json:"no"`
	Link          string    `json:"link"`
	Status        string    `json:"status,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	DigitalAttrs  any       `json:"digital,omitempty"`
	FinanceAttrs  any       `json:"finance"`
	LegalAttrs    any       `json:"legal"`
}

type ApplicationList struct {
	OverallCount uint64         `json:"overall_count"`
	Applications []*Application `json:"applications"`
}

func EncodeApplication(input *domain.Application) *Application {
	if input == nil {
		return nil
	}

	return &Application{
		ID:            input.ID,
		SpreadsheetID: input.SpreadsheetID,
		ApplicationNo: input.ApplicationNo,
		Link:          input.Link,
		Status:        input.Status,
		CreatedAt:     input.CreatedAt,
		DigitalAttrs:  input.DigitalAttrs,
		FinanceAttrs:  input.FinanceAttrs,
		LegalAttrs:    input.LegalAttrs,
	}
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
