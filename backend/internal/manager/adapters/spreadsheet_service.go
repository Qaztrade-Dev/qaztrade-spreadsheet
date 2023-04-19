package adapters

import (
	"context"
	"fmt"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SpreadsheetServiceGoogle struct {
	driveSvc  *drive.Service
	sheetsSvc *sheets.Service
}

var _ domain.SpreadsheetService = (*SpreadsheetServiceGoogle)(nil)

func NewSpreadsheetService(ctx context.Context, credentialsJson []byte) (*SpreadsheetServiceGoogle, error) {
	driveSvc, err := drive.NewService(ctx, option.WithCredentialsJSON(credentialsJson))
	if err != nil {
		return nil, err
	}

	sheetsSvc, err := sheets.NewService(ctx, option.WithCredentialsJSON(credentialsJson))
	if err != nil {
		return nil, err
	}

	return &SpreadsheetServiceGoogle{
		driveSvc:  driveSvc,
		sheetsSvc: sheetsSvc,
	}, err
}

func (s *SpreadsheetServiceGoogle) SwitchModeRead(ctx context.Context, spreadsheetID string) error {
	permission := &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}
	_, err := s.driveSvc.Permissions.Create(spreadsheetID, permission).Do()
	if err != nil {
		return err
	}
	return nil
}

func (s *SpreadsheetServiceGoogle) SwitchModeEdit(ctx context.Context, spreadsheetID string) error {
	permission := &drive.Permission{
		Type: "anyone",
		Role: "writer",
	}
	_, err := s.driveSvc.Permissions.Create(spreadsheetID, permission).Do()
	if err != nil {
		return err
	}
	return nil
}

func (s *SpreadsheetServiceGoogle) Comments(ctx context.Context, spreadsheetID string) error {
	// Get all comments in the sheet
	comments, err := s.driveSvc.Comments.List(spreadsheetID).Fields("*").Do()
	if err != nil {
		return err
	}

	// Print all comments and their replies
	for _, comment := range comments.Comments {
		fmt.Printf("Comment: %s\n", comment.HtmlContent)
		fmt.Printf("Author: %s\n", comment.Author.DisplayName)
		fmt.Printf("Anchor: %s\n", comment.Anchor)
		fmt.Printf("Anchor: %s\n", comment.Kind)
		fmt.Printf("QuotedFileContent: %#v\n", comment.QuotedFileContent)
		for _, reply := range comment.Replies {
			fmt.Printf("\treply: %s\n\tAuthor: %s\n", reply.HtmlContent, reply.Author.DisplayName)
		}

		// resp, err := s.sheetsSvc.Spreadsheets.GetByDataFilter(spreadsheetID, &sheets.GetSpreadsheetByDataFilterRequest{
		// 	DataFilters: []*sheets.DataFilter{
		// 		{
		// 			A1Range: "428107940",
		// 		},
		// 	},
		// }).Do()
		// fmt.Println(resp)
		// fmt.Println(err)
		// if err != nil {
		// 	return err
		// }

		// break
	}

	return nil
}
