package domain

import (
	"context"
	"io"
)

type (
	UpdateCellInput struct {
		SheetID   int64
		RowIdx    int64
		ColumnIdx int64
		Value     string
	}

	SheetsRepository interface {
		UpdateApplication(ctx context.Context, spreadsheetID string, application *Application) error
		UpdateCell(ctx context.Context, spreadsheetID string, input *UpdateCellInput) error
	}

	Storage interface {
		Upload(ctx context.Context, folderName, fileName string, fileSize int64, fileReader io.Reader) (string, error)
		Remove(ctx context.Context, filePath string) error
	}
)

type SpreadsheetClaims struct {
	SpreadsheetID string `json:"sid"`
}
