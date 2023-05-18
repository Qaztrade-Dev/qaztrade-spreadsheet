package domain

import (
	"bytes"
	"errors"
	"io"

	"golang.org/x/net/context"
)

type PDFService interface {
	Create(application *Application, attachments []io.ReadSeeker) (*bytes.Buffer, error)
}

var (
	ErrorSpreadsheetHasMergedCells = errors.New("Таблица содержит объединенные ячейки! ⛔️ Объединенные ячейки запрещены.")
	ErrorAbsentExpenses            = errors.New("Таблица не содержит затраты!")
	ErrorExpensesZero              = errors.New("Заявленные затраты равны нулю! ⛔️ Запрещено подавать заявку на сумму 0 тенге.")
)

type SpreadsheetRepository interface {
	GetApplication(ctx context.Context, spreadsheetID string) (*Application, error)
	GetExpensesSheetTitles(ctx context.Context, spreadsheetID string) ([]string, error)
	GetExpenseValues(ctx context.Context, spreadsheetID string, expensesTitles []string) ([]float64, error)
	GetAttachments(ctx context.Context, spreadsheetID string, expensesTitles []string) ([]io.ReadSeeker, error)
	UpdateSigningTime(ctx context.Context, spreadsheetID, signingTime string) error
	SwitchModeRead(ctx context.Context, spreadsheetID string) error
	HasMergedCells(ctx context.Context, spreadsheetID string, expensesTitles []string) (bool, error)
}

type CreateSigningDocumentResponse struct {
	DocumentID string
	SignLink   string
}

type SigningService interface {
	CreateSigningDocument(ctx context.Context, documentName string, documentReader io.Reader) (*CreateSigningDocumentResponse, error)
}

const (
	StatusUserFilling      = "user_filling"
	StatusManagerReviewing = "manager_reviewing"
	StatusUserFixing       = "user_fixing"
	StatusCompleted        = "completed"
	StatusRejected         = "rejected"
)

type SignApplication struct {
	SpreadsheetID string
	SignLink      string
	Status        string
}

type ApplicationRepository interface {
	AssignSigningInfo(ctx context.Context, spreadsheetID string, info *CreateSigningDocumentResponse) error
	ConfirmSigningInfo(ctx context.Context, spreadsheetID string) error
	GetApplication(ctx context.Context, spreadsheetID string) (*SignApplication, error)
	EditStatus(ctx context.Context, spreadsheetID, statusName string) error
	GetApplicationByDocumentID(ctx context.Context, documentID string) (*SignApplication, error)
}

type SpreadsheetClaims struct {
	SpreadsheetID string `json:"sid"`
}
