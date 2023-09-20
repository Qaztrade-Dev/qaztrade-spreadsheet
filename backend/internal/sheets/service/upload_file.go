package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"github.com/google/uuid"
	"google.golang.org/api/sheets/v4"
)

type UploadFileRequest struct {
	SpreadsheetID string
	SheetID       int64
	SheetName     string
	RowIdx        int64
	ColumnIdx     int64
	FileReader    io.Reader
	FileSize      int64
	FileName      string
}

func (s *service) UploadFile(ctx context.Context, req *UploadFileRequest) error {
	statusApplication, err := s.applicationRepo.GetApplication(ctx, req.SpreadsheetID)
	if err != nil {
		return err
	}

	if statusApplication.Status == domain.StatusUserFixing {

		var (
			filter []*sheets.DataFilter
			key    = fmt.Sprintf("!%s-%d:%d", req.SheetName, req.RowIdx, req.ColumnIdx)
		)

		filter = append(filter, &sheets.DataFilter{
			DeveloperMetadataLookup: &sheets.DeveloperMetadataLookup{
				MetadataKey: key,
			},
		})
		reqMeta := &sheets.SearchDeveloperMetadataRequest{
			DataFilters: filter,
		}

		response, err := s.spreadsheetDevMetadataSvc.Search(req.SpreadsheetID, reqMeta).Do()
		if err != nil {
			return err
		}
		if len(response.MatchedDeveloperMetadata) != 0 {
			folderName := fmt.Sprintf("%s/%s", req.SpreadsheetID, req.SheetName)
			filekey := fmt.Sprintf("%s/%s-%s", folderName, uuid.NewString(), req.FileName)

			value, err := s.storage.Upload(ctx, filekey, req.FileSize, req.FileReader)
			if err != nil {
				log.Printf("storage.Upload error file: folderName - %s, fileName - %s\n", folderName, req.FileName)
				return err
			}

			if err := s.sheetsRepo.UpdateCell(ctx, req.SpreadsheetID, &domain.UpdateCellInput{
				SheetID:   req.SheetID,
				RowIdx:    req.RowIdx,
				ColumnIdx: req.ColumnIdx,
				Value:     value,
				Replace:   false,
				SheetName: req.SheetName,
			}); err != nil {
				return err
			}
		}

	} else if statusApplication.Status == domain.StatusUserFilling {
		// 2. if the cell contains url, delete the file
		Hyperlink, err := s.sheetsRepo.GetHyperLink(ctx, req.SpreadsheetID, req.SheetName, req.RowIdx, req.ColumnIdx)
		if err != nil {
			return err
		}
		if *Hyperlink != "" {
			filePath := getFilePath(req.SpreadsheetID, *Hyperlink)
			if err := s.storage.Remove(ctx, filePath); err != nil {
				return err
			}
		}

		// 3. upload file, get url
		folderName := fmt.Sprintf("%s/%s", req.SpreadsheetID, req.SheetName)
		filekey := fmt.Sprintf("%s/%s-%s", folderName, uuid.NewString(), req.FileName)

		value, err := s.storage.Upload(ctx, filekey, req.FileSize, req.FileReader)
		if err != nil {
			log.Printf("storage.Upload error file: folderName - %s, fileName - %s\n", folderName, req.FileName)
			return err
		}

		// 4. write url to cell
		if err := s.sheetsRepo.UpdateCell(ctx, req.SpreadsheetID, &domain.UpdateCellInput{
			SheetID:   req.SheetID,
			RowIdx:    req.RowIdx,
			ColumnIdx: req.ColumnIdx,
			Value:     value,
			Replace:   true,
			SheetName: req.SheetName,
		}); err != nil {
			return err
		}
	}
	return nil
}

func getFilePath(spreadsheetID, hyperlink string) string {

	splittedArgs := strings.Split(hyperlink, ";")
	if len(splittedArgs) == 0 {
		return ""
	}

	var (
		quotedLink       = splittedArgs[0]
		spreadsheetIDIdx = strings.Index(quotedLink, spreadsheetID)
		filePath         = quotedLink[spreadsheetIDIdx:]
	)

	return filePath
}
