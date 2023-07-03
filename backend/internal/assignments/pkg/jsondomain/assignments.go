package jsondomain

import (
	"time"

	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
)

type AssignmentView struct {
	ID             int       `json:"id"`
	ApplicantName  string    `json:"applicant_name"`
	ApplicantBIN   string    `json:"applicant_bin"`
	SheetTitle     string    `json:"sheet_title"`
	SheetID        uint64    `json:"sheet_id"`
	AssignmentType string    `json:"assignment_type"`
	Link           string    `json:"link"`
	AssigneeName   string    `json:"assignee_name"`
	TotalRows      int       `json:"total_rows"`
	TotalSum       int       `json:"total_sum"`
	RowsCompleted  int       `json:"rows_completed"`
	IsCompleted    bool      `json:"is_completed"`
	CompletedAt    time.Time `json:"completed_at"`
}

func EncodeAssignmentView(input *domain.AssignmentView) *AssignmentView {
	if input == nil {
		return nil
	}

	return &AssignmentView{
		ID:             input.ID,
		ApplicantName:  input.ApplicantName,
		ApplicantBIN:   input.ApplicantBIN,
		SheetTitle:     input.SheetTitle,
		SheetID:        input.SheetID,
		AssignmentType: input.AssignmentType,
		Link:           input.Link,
		AssigneeName:   input.AssigneeName,
		TotalRows:      input.TotalRows,
		TotalSum:       input.TotalSum,
		RowsCompleted:  input.RowsCompleted,
		IsCompleted:    input.IsCompleted,
		CompletedAt:    input.CompletedAt,
	}
}

type AssignmentsList struct {
	Total   int               `json:"total"`
	Objects []*AssignmentView `json:"objects"`
}

func EncodeAssignmentsList(input *domain.AssignmentsList) *AssignmentsList {
	if input == nil {
		return nil
	}

	return &AssignmentsList{
		Total:   input.Total,
		Objects: EncodeSlice(input.Objects, EncodeAssignmentView),
	}
}

type AssignmentsInfo struct {
	Total     uint64 `json:"total"`
	Completed uint64 `json:"completed"`
}

func EncodeAssignmentsInfo(input *domain.AssignmentsInfo) *AssignmentsInfo {
	if input == nil {
		return nil
	}

	return &AssignmentsInfo{
		Total:     input.Total,
		Completed: input.Completed,
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
