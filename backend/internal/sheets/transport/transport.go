package transport

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/sheets/endpoint"
	"github.com/doodocs/qaztrade/backend/internal/sheets/pkg/tally"
)

func DecodeSubmitApplicationRequest(_ context.Context, r *http.Request) (interface{}, error) {
	tallyJsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	application, err := tally.Decode(tallyJsonBytes)
	if err != nil {
		return nil, err
	}

	return endpoint.SubmitApplicationRequest{
		Application: application,
	}, nil
}

func DecodeUploadFileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return nil, err
	}

	fileReader, header, err := r.FormFile("fileInput")
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	var (
		fileName = header.Filename
		fileSize = header.Size
		jsonData = r.FormValue("selected_cell")
	)

	var body struct {
		SheetName string `json:"sheet_name"`
		SheetID   int64  `json:"sheet_id"`
		RowIdx    int64  `json:"row_idx"`
		ColumnIdx int64  `json:"column_idx"`
	}

	if err := json.Unmarshal([]byte(jsonData), &body); err != nil {
		return nil, err
	}

	var (
		tokenString = extractHeaderToken(r)
	)

	return endpoint.UploadFileRequest{
		SpreadsheetTokenStr: tokenString,
		SheetID:             body.SheetID,
		SheetName:           body.SheetName,
		RowIdx:              body.RowIdx,
		ColumnIdx:           body.ColumnIdx,
		FileReader:          fileReader,
		FileSize:            fileSize,
		FileName:            fileName,
	}, nil
}

func extractHeaderToken(r *http.Request) string {
	authorization := r.Header.Get("authorization")
	if authorization == "" {
		return ""
	}

	tokenString := strings.Split(authorization, " ")[1]
	return tokenString
}

func extractHeaderSheetID(r *http.Request) (int64, error) {
	sheetIDStr := r.Header.Get("x-sheet-id")
	sheetID, err := strconv.ParseInt(sheetIDStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return sheetID, nil
}

func extractHeaderSheetName(r *http.Request) string {
	sheetName := r.Header.Get("x-sheet-name")
	unescaped, _ := url.QueryUnescape(sheetName)
	return unescaped
}
