package service

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
)

type UploadFileRequest struct {
	SpreadsheetID string
	SheetID       int64
	SheetName     string
	Hyperlink     string
	RowIdx        int64
	ColumnIdx     int64

	FileReader io.Reader
	FileSize   int64
	FileName   string
}

func (s *service) UploadFile(ctx context.Context, req *UploadFileRequest) error {
	// 2. if the cell contains url, delete the file
	if req.Hyperlink != "" {
		filePath := getFilePath(req.SpreadsheetID, req.Hyperlink)
		if err := s.storage.Remove(ctx, filePath); err != nil {
			return err
		}
	}

	// 3. upload file, get url
	folderName := fmt.Sprintf("%s/%s", req.SpreadsheetID, req.SheetName)

	value, err := s.storage.Upload(ctx, folderName, req.FileName, req.FileSize, req.FileReader)
	if err != nil {
		return err
	}

	valueHyperlink := fmt.Sprintf("=HYPERLINK(\"%s\"; \"файл\")", value)

	// 4. write url to cell
	if err := s.sheetsRepo.UpdateCell(ctx, req.SpreadsheetID, &domain.UpdateCellInput{
		SheetID:   req.SheetID,
		RowIdx:    req.RowIdx,
		ColumnIdx: req.ColumnIdx,
		Value:     valueHyperlink,
	}); err != nil {
		return err
	}

	return nil
}

func getFilePath(spreadsheetID, hyperlink string) string {
	afterHyperlink, ok := strings.CutPrefix(hyperlink, "=HYPERLINK(")
	if !ok {
		return ""
	}

	splittedArgs := strings.Split(afterHyperlink, ";")
	if len(splittedArgs) == 0 {
		return ""
	}

	var (
		quotedLink       = splittedArgs[0]
		link             = strings.ReplaceAll(quotedLink, "\"", "")
		spreadsheetIDIdx = strings.Index(link, spreadsheetID)
		filePath         = link[spreadsheetIDIdx:]
	)

	return filePath
}
