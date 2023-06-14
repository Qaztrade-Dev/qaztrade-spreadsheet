package domain

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"golang.org/x/net/context"
)

type PDFService interface {
	Create(application *Application, attachments []io.ReadSeeker) (*bytes.Buffer, error)
}

var (
	ErrorEmptySpreadsheet          = errors.New("Таблица не содержит видов затрат!")
	ErrorSpreadsheetHasMergedCells = errors.New("Таблица содержит объединенные ячейки! ⛔️ Объединенные ячейки запрещены.")
	ErrorAbsentExpenses            = errors.New("Таблица не содержит затраты!")
	ErrorExpensesZero              = errors.New("Заявленные затраты равны нулю! ⛔️ Запрещено подавать заявку на сумму 0 тенге.")
)

type Sheet struct {
	Title    string
	SheetID  int64
	Expenses float64
	Rows     int64
}

func SheetTitles(input []*Sheet) []string {
	titles := make([]string, 0, len(input))
	for _, sheet := range input {
		title := sheet.Title
		titles = append(titles, title)
	}
	return titles
}

func SheetTitlesJoined(input []*Sheet) string {
	titles := SheetTitles(input)
	return strings.Join(titles, ", ")
}

func SheetsSumIsZero(input []*Sheet) bool {
	sum := SheetsTotalExpenses(input)
	return sum == 0
}

func SheetsTotalExpenses(input []*Sheet) float64 {
	sum := float64(0)
	for _, sheet := range input {
		sum += sheet.Expenses
	}
	return sum
}

func SheetsTotalRows(input []*Sheet) int64 {
	sum := int64(0)
	for _, sheet := range input {
		sum += sheet.Rows
	}
	return sum
}

func SheetsAgg(input []*Sheet) *Sheet {
	return &Sheet{
		Expenses: SheetsTotalExpenses(input),
		Rows:     SheetsTotalRows(input),
	}
}

func GetApplicationDate() string {
	now := time.Now()

	location, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		return now.Format("02.01.2006")
	}

	timeStr := now.In(location).Format("02.01.2006")
	return timeStr
}

func GetDocumentName(bin string) string {
	var (
		timeStr      = GetApplicationDate()
		documentName = fmt.Sprintf("Заявление %s %s", bin, timeStr)
	)
	return documentName
}

type ApplicationAttrs struct {
	Application *Application
	SheetsAgg   *Sheet
	Sheets      []*Sheet
}

type SpreadsheetRepository interface {
	GetApplication(ctx context.Context, spreadsheetID string) (*Application, error)

	// GetSheets returns SheetInformation for each sheet in a spreadsheet.
	GetSheets(ctx context.Context, spreadsheetID string) ([]*Sheet, error)

	GetAttachments(ctx context.Context, spreadsheetID string, sheets []*Sheet) ([]io.ReadSeeker, error)
	UpdateSigningTime(ctx context.Context, spreadsheetID, signingTime string) error
	SwitchModeRead(ctx context.Context, spreadsheetID string) error
	BlockImportantRanges(ctx context.Context, spreadsheetID string) error
	HasMergedCells(ctx context.Context, spreadsheetID string, sheets []*Sheet) (bool, error)
}

type CreateSigningDocumentResponse struct {
	DocumentID string
	SignLink   string
}

type SigningService interface {
	GetSigningTime(ctx context.Context, documentID string) (time.Time, error)
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
	AssignAttrs(ctx context.Context, spreadsheetID string, input *ApplicationAttrs) error
	ConfirmSigningInfo(ctx context.Context, spreadsheetID, signedAtStr string) error
	GetApplication(ctx context.Context, spreadsheetID string) (*SignApplication, error)
	EditStatus(ctx context.Context, spreadsheetID, statusName string) error
	GetApplicationByDocumentID(ctx context.Context, documentID string) (*SignApplication, error)
}

type SpreadsheetClaims struct {
	SpreadsheetID string `json:"sid"`
}
