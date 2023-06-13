package adapters

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sign/domain"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ApplicationRepositoryPostgre struct {
	pg *pgxpool.Pool
}

var _ domain.ApplicationRepository = (*ApplicationRepositoryPostgre)(nil)

func NewApplicationRepositoryPostgre(pg *pgxpool.Pool) *ApplicationRepositoryPostgre {
	return &ApplicationRepositoryPostgre{
		pg: pg,
	}
}

func (r *ApplicationRepositoryPostgre) AssignSigningInfo(ctx context.Context, spreadsheetID string, info *domain.CreateSigningDocumentResponse) error {
	const sql = `
		update "applications" set
			sign_link=$1,
			sign_document_id=$2
		where 
			spreadsheet_id=$3
	`

	if _, err := r.pg.Exec(ctx, sql, info.SignLink, info.DocumentID, spreadsheetID); err != nil {
		return err
	}

	return nil
}

func (r *ApplicationRepositoryPostgre) ConfirmSigningInfo(ctx context.Context, spreadsheetID, signedAtStr string) error {
	const sql = `
		update "applications" set
			is_signed=true,
			sign_at=TO_TIMESTAMP($2, 'DD.MM.YYYY')
		where 
			spreadsheet_id=$1
	`

	if _, err := r.pg.Exec(ctx, sql, spreadsheetID, signedAtStr); err != nil {
		return err
	}

	return nil
}

func (r *ApplicationRepositoryPostgre) EditStatus(ctx context.Context, spreadsheetID, statusName string) error {
	const sql = `
		update "applications" set
			status_id = (select id from "application_statuses" where value = $2)
		where 
			spreadsheet_id = $1
	`

	if _, err := r.pg.Exec(ctx, sql, spreadsheetID, statusName); err != nil {
		return err
	}
	return nil
}

func (r *ApplicationRepositoryPostgre) GetApplication(ctx context.Context, spreadsheetID string) (*domain.SignApplication, error) {
	const query = `
		select
			a.sign_link,
			aps.value
		from "applications" a
		join "application_statuses" aps on aps.id = a.status_id
		where 
			a.spreadsheet_id = $1
	`

	var (
		scanSignLink *string
		scanStatus   *string
	)

	if err := r.pg.QueryRow(ctx, query, spreadsheetID).Scan(
		&scanSignLink,
		&scanStatus,
	); err != nil {
		return nil, err
	}

	result := &domain.SignApplication{
		SignLink: valueFromPointer(scanSignLink),
		Status:   valueFromPointer(scanStatus),
	}

	return result, nil
}

func (r *ApplicationRepositoryPostgre) GetApplicationByDocumentID(ctx context.Context, documentID string) (*domain.SignApplication, error) {
	const query = `
		select
			a.spreadsheet_id,
			a.sign_link,
			aps.value
		from "applications" a
		join "application_statuses" aps on aps.id = a.status_id
		where 
			a.sign_document_id = $1
	`

	var (
		scanSpreadsheetID *string
		scanSignLink      *string
		scanStatus        *string
	)

	if err := r.pg.QueryRow(ctx, query, documentID).Scan(
		&scanSpreadsheetID,
		&scanSignLink,
		&scanStatus,
	); err != nil {
		return nil, err
	}

	result := &domain.SignApplication{
		SpreadsheetID: valueFromPointer(scanSpreadsheetID),
		SignLink:      valueFromPointer(scanSignLink),
		Status:        valueFromPointer(scanStatus),
	}

	return result, nil
}

func valueFromPointer[T any](value *T) T {
	var defaultValue T

	if value == nil {
		return defaultValue
	}
	return *value
}
