package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
)

type SubmitRecordRequest struct {
	SpreadsheetID string
	SheetName     string
	SheetID       int64
	Payload       *domain.Payload
}

func (s *service) SubmitRecord(ctx context.Context, req *SubmitRecordRequest) error {
	if err := s.traversePayload(ctx, req.Payload.Value); err != nil {
		return err
	}

	if err := s.sheetsRepo.InsertRecord(ctx, req.SpreadsheetID, req.SheetName, req.SheetID, req.Payload); err != nil {
		return err
	}

	return nil
}

func (s *service) traversePayload(ctx context.Context, payload map[string]interface{}) error {
	for k := range payload {
		value := payload[k]
		switch {
		case isFile(value):
			file := decodeFile(value)
			if file.FileSize == 0 {
				delete(payload, k)
				continue
			}
			value, err := s.storage.Upload(ctx, "folder", file.FileName, file.FileSize, bytes.NewReader(file.File))
			if err != nil {
				return err
			}
			payload[k] = fmt.Sprintf("=HYPERLINK(\"%s\", \"файл\")", value)
		case isPayload(value):
			return s.traversePayload(ctx, value.(map[string]interface{}))
		}
	}
	return nil
}

func isFile(value interface{}) bool {
	mapValue, ok := value.(map[string]interface{})
	if !ok {
		return false
	}

	if _, ok := mapValue["file"]; !ok {
		return false
	}

	return true
}

type FilePayload struct {
	File     []byte `json:"file"`
	FileSize int64  `json:"size"`
	FileName string `json:"name"`
}

func decodeFile(value interface{}) *FilePayload {
	var (
		mapValue = value.(map[string]interface{})
		fileB64  = mapValue["file"].(string)
		fileSize = mapValue["size"].(float64)
		fileName = mapValue["name"].(string)

		findSubstr = "base64,"
		substrIdx  = strings.Index(fileB64, findSubstr)
	)

	fileBytes, _ := base64.StdEncoding.DecodeString(fileB64[substrIdx+len(findSubstr):])

	return &FilePayload{
		File:     fileBytes,
		FileName: fileName,
		FileSize: int64(fileSize),
	}
}

func isPayload(value interface{}) bool {
	_, ok := value.(map[string]interface{})
	return ok
}
