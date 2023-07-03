package adapters

import (
	"context"
	"encoding/json"
	"time"

	"github.com/doodocs/qaztrade/backend/internal/sign/domain"
	"github.com/doodocs/qaztrade/backend/internal/sign/pkg/jsondomain"
	"github.com/doodocs/qaztrade/backend/pkg/postgres"
	"github.com/jackc/pgx/v4"
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

func (r *ApplicationRepositoryPostgre) AssignSigningInfo(ctx context.Context, spreadsheetID string, input *domain.CreateSigningDocumentResponse) error {
	const sql = `
		update "applications" set
			sign_link=$1,
			sign_document_id=$2
		where 
			spreadsheet_id=$3
	`

	if _, err := r.pg.Exec(ctx, sql, input.SignLink, input.DocumentID, spreadsheetID); err != nil {
		return err
	}

	return nil
}

func (r *ApplicationRepositoryPostgre) AssignAttrs(ctx context.Context, spreadsheetID string, input *domain.ApplicationAttrs) error {
	const sql = `
		update "applications" set
			attrs=$2
		where 
			spreadsheet_id=$1
	`

	attrs := jsondomain.EncodeApplicationAttrs(input)
	attrsBytes, err := json.Marshal(attrs)
	if err != nil {
		return err
	}

	if _, err := r.pg.Exec(ctx, sql, spreadsheetID, string(attrsBytes)); err != nil {
		return err
	}

	return nil
}

func (r *ApplicationRepositoryPostgre) ConfirmSigningInfo(ctx context.Context, spreadsheetID string, signedAt time.Time) error {
	err := postgres.InTransaction(ctx, r.pg, func(ctx context.Context, tx pgx.Tx) error {
		no, err := assignApplicationNo(ctx, tx, spreadsheetID)
		if err != nil {
			return nil
		}

		if err := confirmApplicationSigning(ctx, tx, spreadsheetID, signedAt, no); err != nil {
			return nil
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func confirmApplicationSigning(ctx context.Context, tx pgx.Tx, spreadsheetID string, signedAt time.Time, no int) error {
	signedAtStr := signedAt.Format(domain.TimestampLayout)
	const sql = `
		update "applications" set
			is_signed=true,
			sign_at=$2,
			no=$3
		where 
			spreadsheet_id=$1
	`

	if _, err := tx.Exec(ctx, sql, spreadsheetID, signedAtStr, no); err != nil {
		return err
	}

	return nil
}

func assignApplicationNo(ctx context.Context, tx pgx.Tx, spreadsheetID string) (int, error) {
	const sql = `
		insert into "application_signings" (application_id)
		select
			id
		from applications
		where
			spreadsheet_id=$1
		returning id
	`

	var no int
	if err := tx.QueryRow(ctx, sql, spreadsheetID).Scan(&no); err != nil {
		return -1, err
	}

	return no, nil
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
		SignLink: postgres.Value(scanSignLink),
		Status:   postgres.Value(scanStatus),
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
		SpreadsheetID: postgres.Value(scanSpreadsheetID),
		SignLink:      postgres.Value(scanSignLink),
		Status:        postgres.Value(scanStatus),
	}

	return result, nil
}
