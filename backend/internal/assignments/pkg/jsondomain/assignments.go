package jsondomain

import (
	"fmt"
	"time"

	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
)

type AssignmentView struct {
	ApplicationID    string     `json:"application_id"`
	AssignmentID     uint64     `json:"assignment_id"`
	ID               int64      `json:"id"`
	ApplicantName    string     `json:"applicant_name"`
	ApplicantBIN     string     `json:"applicant_bin"`
	SheetTitle       string     `json:"sheet_title"`
	SheetID          uint64     `json:"sheet_id"`
	AssignmentType   string     `json:"assignment_type"`
	Link             string     `json:"link"`
	SignLink         string     `json:"sign_link"`
	AssigneeName     string     `json:"assignee_name"`
	TotalRows        int        `json:"total_rows"`
	TotalSum         int        `json:"total_sum"`
	RowsCompleted    int        `json:"rows_completed"`
	IsCompleted      bool       `json:"is_completed"`
	CompletedAt      *time.Time `json:"completed_at"`
	ResolutionStatus string     `json:"resolution_status"`
	ResolvedAt       *time.Time `json:"resolved_at"`
	ReplyEndAt       *time.Time `json:"reply_end_at"`
	DigitalStatus    string     `json:"digital_status,omitempty"`
	FinanceStatus    string     `json:"finance_status"`
	LegalStatus      string     `json:"legal_status"`
}

func EncodeAssignmentView(input *domain.AssignmentView) *AssignmentView {
	if input == nil {
		return nil
	}

	var replyEndAt time.Time
	if !input.ResolvedAt.IsZero() {
		replyEndAt = input.ResolvedAt.Add(input.CountdownDuration)
	}

	return &AssignmentView{
		ApplicationID:    input.ApplicationID,
		AssignmentID:     input.AssignmentID,
		ID:               input.ID,
		ApplicantName:    input.ApplicantName,
		ApplicantBIN:     input.ApplicantBIN,
		SheetTitle:       input.SheetTitle,
		SheetID:          input.SheetID,
		AssignmentType:   input.AssignmentType,
		Link:             input.Link,
		SignLink:         input.SignLink,
		AssigneeName:     input.AssigneeName,
		TotalRows:        input.TotalRows,
		TotalSum:         input.TotalSum,
		RowsCompleted:    input.RowsCompleted,
		IsCompleted:      input.IsCompleted,
		CompletedAt:      timeToPtr(input.CompletedAt),
		ResolutionStatus: input.ResolutionStatus,
		ResolvedAt:       timeToPtr(input.ResolvedAt),
		ReplyEndAt:       timeToPtr(replyEndAt),
		DigitalStatus:    input.DigitalStatus,
		FinanceStatus:    input.FinanceStatus,
		LegalStatus:      input.LegalStatus,
	}
}

func timeToPtr(input time.Time) *time.Time {
	if input.IsZero() {
		return nil
	}
	return &input
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

type DialogMessage struct {
	SentAt          *time.Time `json:"sent_at"`
	UserDisplayName *string    `json:"user_display_name"`
	UserEmail       string     `json:"user_email"`
	UserID          string     `json:"user_id"`
	Body            string     `json:"body"`
}

func EncodeDialogMessage(input *domain.Message) *DialogMessage {
	if input == nil {
		return nil
	}

	msg := &DialogMessage{
		SentAt:    timeToPtr(input.CreatedAt),
		UserEmail: input.Email,
		UserID:    input.UserID,
	}

	if input.FullName != "" {
		msg.UserDisplayName = &input.FullName
	}

	var body string

	if _, ok := input.Attrs["file_url"]; ok {
		body = dialogMsgManagerNotice(input)
	}

	if _, ok := input.Attrs["sign_link"]; ok {
		body = dialogMsgUserRespond(input)
	}

	msg.Body = body

	return msg
}

func dialogMsgManagerNotice(input *domain.Message) string {
	return fmt.Sprintf(`Отправлен [Документ-уведомеление](%s)`, input.Attrs["file_url"])
}

func dialogMsgUserRespond(input *domain.Message) string {
	var body string

	if input.DoodocsIsSigned {
		body = fmt.Sprintf(`Отправлено [Сопроводительное письмо](%s)`, input.Attrs["sign_link"])
		body += "\n\n✅ Подписано"
		body += fmt.Sprintf("\n\n✅ Время подписания: %s", input.CreatedAt.Format("02.01.2006 15:04:05"))
	} else {
		body = fmt.Sprintf(`Создано [Сопроводительное письмо](%s)`, input.Attrs["sign_link"])
		body += "\n\n⏳ Ожидается подписание"
	}

	return body
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
