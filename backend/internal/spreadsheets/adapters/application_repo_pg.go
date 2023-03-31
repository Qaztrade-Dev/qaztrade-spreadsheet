package adapters

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ApplicationRepositoryPostgre struct {
	pg *pgxpool.Pool
}

var _ domain.ApplicationRepository = (*ApplicationRepositoryPostgre)(nil)

func NewAuthorizationRepositoryPostgre(pg *pgxpool.Pool) *ApplicationRepositoryPostgre {
	return &ApplicationRepositoryPostgre{
		pg: pg,
	}
}

func (r *ApplicationRepositoryPostgre) Create(ctx context.Context, userID string, input *domain.Application) error {
	const sql = `
		insert into "users" 
			(
				user_id,
				status_id,
				spreadsheet_id,
				link
			)
		values
			(
				$1,
				(select id from "application_statuses" where value = 'user_filling'),
				$2,
				$3
			)
	`

	if _, err := r.pg.Exec(ctx, sql, userID, input.SpreadsheetID, input.Link); err != nil {
		return err
	}

	return nil
}
