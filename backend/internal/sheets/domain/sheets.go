package domain

import (
	"context"
	"io"
)

type (
	UpdateCellInput struct {
		SheetID       int64
		PrevHyperlink string
		SheetName     string
		RowIdx        int64
		ColumnIdx     int64
		Value         string
		Replace       bool
	}

	AddRowsInput struct {
		SheetID    int64
		SheetName  string
		RowsAmount int
	}

	SheetsRepository interface {
		UpdateApplication(ctx context.Context, spreadsheetID string, application *Application) error
		UpdateCell(ctx context.Context, spreadsheetID string, input *UpdateCellInput) error
		AddRows(ctx context.Context, spreadsheetID string, input *AddRowsInput) error
		GetHyperLink(ctx context.Context, spreadsheetID string, SheetName string, Row_idx int64, Column_idx int64) (*string, error)
	}

	Storage interface {
		Upload(ctx context.Context, filekey string, fileSize int64, fileReader io.Reader) (string, error)
		Remove(ctx context.Context, filePath string) error
	}
)

type SpreadsheetClaims struct {
	SpreadsheetID string `json:"sid"`
}

type ApplicationRepository interface {
	GetApplication(ctx context.Context, spreadsheetID string) (*StatusApplication, error)
}

type SpreadsheetService interface {
}
