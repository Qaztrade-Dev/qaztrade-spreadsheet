package adapters

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"github.com/doodocs/qaztrade/backend/pkg/postgres"
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

func (r *ApplicationRepositoryPostgre) GetApplication(ctx context.Context, spreadsheetID string) (*domain.StatusApplication, error) {
	const query = `
		select
			aps.value
		from "applications" a
		join "application_statuses" aps on aps.id = a.status_id
		where 
			a.spreadsheet_id = $1
	`

	var (
		scanStatus *string
	)

	if err := r.pg.QueryRow(ctx, query, spreadsheetID).Scan(
		&scanStatus,
	); err != nil {
		return nil, err
	}

	result := &domain.StatusApplication{
		Status: postgres.Value(scanStatus),
	}

	return result, nil
}
