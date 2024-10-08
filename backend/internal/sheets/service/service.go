package service

import (
	"context"

	assignmentsDomain "github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"google.golang.org/api/sheets/v4"
)

type Service interface {
	SubmitApplication(ctx context.Context, req *SubmitApplicationRequest) error
	UploadFile(ctx context.Context, req *UploadFileRequest) error
	AddRows(ctx context.Context, req *AddRowsRequest) error
}

type service struct {
	sheetsRepo                domain.SheetsRepository
	storage                   domain.Storage
	applicationRepo           domain.ApplicationRepository
	spreadsheetDevMetadataSvc sheets.SpreadsheetsDeveloperMetadataService
	assignmentsRepo           assignmentsDomain.AssignmentsRepository
}

func NewService(
	sheetsRepo domain.SheetsRepository,
	storage domain.Storage,
	applicationRepo domain.ApplicationRepository,
	spreadsheetDevMetadataSvc sheets.SpreadsheetsDeveloperMetadataService,
	assignmentsRepo assignmentsDomain.AssignmentsRepository,
) Service {
	return &service{
		sheetsRepo:                sheetsRepo,
		storage:                   storage,
		applicationRepo:           applicationRepo,
		spreadsheetDevMetadataSvc: spreadsheetDevMetadataSvc,
		assignmentsRepo:           assignmentsRepo,
	}
}
