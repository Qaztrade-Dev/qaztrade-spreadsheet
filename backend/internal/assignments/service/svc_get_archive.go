package service

import (
	"context"
	"fmt"
	"io"

	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"github.com/doodocs/qaztrade/backend/pkg/storage"
)

type GetArchiveRequest struct {
	AssignmentID uint64
}

type GetArchiveResponse struct {
	ArchiveReader io.ReadCloser
	RemoveFunc    storage.RemoveFunction
}

func (s *service) GetArchive(ctx context.Context, req *GetArchiveRequest) (*GetArchiveResponse, error) {
	assignment, err := s.assignmentRepo.GetOne(ctx, &domain.GetManyInput{
		AssignmentID: &req.AssignmentID,
	})
	if err != nil {
		return nil, err
	}

	folderName := fmt.Sprintf("%s/%s", assignment.SpreadsheetID, assignment.SheetTitle)

	readCloser, removeFunc, err := s.storage.GetArchive(ctx, folderName)
	if err != nil {
		return nil, err
	}

	resp := &GetArchiveResponse{
		ArchiveReader: readCloser,
		RemoveFunc:    removeFunc,
	}

	return resp, nil
}
