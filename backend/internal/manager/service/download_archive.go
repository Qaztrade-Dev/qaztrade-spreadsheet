package service

import (
	"context"
	"io"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
)

type DownloadArchiveRequest struct {
	ApplicationID string
}

type DownloadArchiveResponse struct {
	ArchiveReader io.ReadCloser
	RemoveFunc    domain.RemoveFunction
}

func (s *service) DownloadArchive(ctx context.Context, req *DownloadArchiveRequest) (*DownloadArchiveResponse, error) {
	application, err := s.applicationRepo.GetOne(ctx, &domain.ApplicationQuery{
		ApplicationID: req.ApplicationID,
	})
	if err != nil {
		return nil, err
	}

	readCloser, removeFunc, err := s.spreadsheetStorage.DownloadArchive(ctx, application.SpreadsheetID)
	if err != nil {
		return nil, err
	}

	resp := &DownloadArchiveResponse{
		ArchiveReader: readCloser,
		RemoveFunc:    removeFunc,
	}

	return resp, nil
}
