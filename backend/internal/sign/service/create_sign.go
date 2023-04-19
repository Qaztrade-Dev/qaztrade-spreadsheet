package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/doodocs/qaztrade/backend/internal/sign/domain"
)

type CreateSignRequest struct {
	SpreadsheetID string
}

func (s *service) CreateSign(ctx context.Context, req *CreateSignRequest) (string, error) {
	var linkbase = "https://link.doodocs.kz/"

	signApplication, err := s.applicationRepo.GetApplication(ctx, req.SpreadsheetID)
	if err != nil {
		return "", err
	}

	if signApplication.Status != domain.StatusUserFilling {
		return linkbase + signApplication.SignLink, nil
	}

	application, err := s.spreadsheetRepo.GetApplication(ctx, req.SpreadsheetID)
	if err != nil {
		return "", err
	}

	expensesTitles, err := s.spreadsheetRepo.GetExpensesSheetTitles(ctx, req.SpreadsheetID)
	if err != nil {
		return "", err
	}

	if len(expensesTitles) == 0 {
		return "", errors.New("no expenses")
	}

	expensesValues, err := s.spreadsheetRepo.GetExpenseValues(ctx, req.SpreadsheetID, expensesTitles)
	if err != nil {
		return "", err
	}

	application.ExpensesList = strings.Join(expensesTitles, ", ")
	application.ExpensesSum = fmt.Sprintf("%v", sumFloats64(expensesValues))
	application.ApplicationDate = createApplicationDate()

	attachments, err := s.spreadsheetRepo.GetAttachments(ctx, req.SpreadsheetID, expensesTitles)
	if err != nil {
		return "", err
	}

	pdfToSign, err := s.pdfSvc.Create(application, attachments)
	if err != nil {
		return "", err
	}

	documentName, err := createDocumentName(application)
	if err != nil {
		return "", err
	}

	resp, err := s.signSvc.CreateSigningDocument(ctx, documentName, pdfToSign)
	if err != nil {
		return "", err
	}

	if err := s.applicationRepo.AssignSigningInfo(ctx, req.SpreadsheetID, resp); err != nil {
		return "", err
	}

	return linkbase + resp.SignLink, nil
}

func createDocumentName(application *domain.Application) (string, error) {
	now := time.Now()

	location, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		return "", err
	}

	timeStr := now.In(location).Format(time.DateTime)
	return fmt.Sprintf("Заявление %s %s", application.Bin, timeStr), nil
}

func createApplicationDate() string {
	now := time.Now()

	location, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		return now.Format("02.01.2006")
	}

	timeStr := now.In(location).Format("02.01.2006")
	return timeStr
}

func sumFloats64(ints []float64) float64 {
	result := float64(0)
	for _, v := range ints {
		result += v
	}
	return result
}
