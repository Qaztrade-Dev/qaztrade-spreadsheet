package adapters

import (
	"context"
	"time"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
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

func (r *ApplicationRepositoryPostgre) Create(ctx context.Context, userID string, input *domain.Application) error {
	const sql = `
		insert into "applications" 
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

func (r *ApplicationRepositoryPostgre) GetMany(ctx context.Context, query *domain.ApplicationQuery) (*domain.ApplicationList, error) {
	applicationsCount, err := r.getCount(ctx, query)
	if err != nil {
		return nil, err
	}

	applications, err := r.getMany(ctx, query)
	if err != nil {
		return nil, err
	}

	result := &domain.ApplicationList{
		OverallCount: applicationsCount,
		Applications: applications,
	}

	return result, nil
}

func (r *ApplicationRepositoryPostgre) getMany(ctx context.Context, query *domain.ApplicationQuery) ([]*domain.Application, error) {
	const sql = `
		with "application_progress" as (
			SELECT distinct on (ass.application_id)
				ass.application_id,
				CASE 
					WHEN digital.id IS NULL THEN NULL
					ELSE jsonb_build_object(
						'user_id', digital.user_id,
						'status', COALESCE(digital_status.value, 'manager_reviewing'),
						'resolved_at', digital.resolved_at,
						'reply_end_at', digital.resolved_at + digital.countdown_duration
					)
				END as digital, 
				CASE 
					WHEN finance.id IS NULL THEN NULL
					ELSE jsonb_build_object(
						'user_id', finance.user_id,
						'status', COALESCE(finance_status.value, 'manager_reviewing'),
						'resolved_at', finance.resolved_at,
						'reply_end_at', finance.resolved_at + finance.countdown_duration
					)
				END as finance, 
				CASE 
					WHEN legal.id IS NULL THEN NULL
					ELSE jsonb_build_object(
						'user_id', legal.user_id,
						'status', COALESCE(legal_status.value, 'manager_reviewing'),
						'resolved_at', legal.resolved_at,
						'reply_end_at', legal.resolved_at + legal.countdown_duration
					)
				END as legal 
			from assignments ass
			left join assignments digital on digital.application_id = ass.application_id and digital.type = 'digital'
			left join application_statuses digital_status on digital_status.id = digital.resolution_status_id
			left join assignments finance on finance.application_id = ass.application_id and finance.type = 'finance'
			left join application_statuses finance_status on finance_status.id = finance.resolution_status_id
			left join assignments legal on legal.application_id = ass.application_id and legal.type = 'legal'
			left join application_statuses legal_status on legal_status.id = legal.resolution_status_id
		)
		select 
			a.created_at,
			ast.value,
			a.spreadsheet_id,
			a.link,
			a.no,
			apr.digital,
			apr.finance,
			apr.legal
		from "applications" a
		join "application_statuses" ast on ast.id = a.status_id
		left join "application_progress" apr on apr.application_id = a.id
		where 
			a.user_id = $1
		order by a.created_at desc
		limit $2
		offset $3
	`

	applications, err := queryApplications(ctx, r.pg, sql, query.UserID, query.Limit, query.Offset)
	if err != nil {
		return nil, err
	}

	return applications, nil
}

func (r *ApplicationRepositoryPostgre) getCount(ctx context.Context, query *domain.ApplicationQuery) (uint64, error) {
	const sql = `
		select 
			count(*)
		from "applications"
		where 
			user_id = $1
	`

	var count uint64
	if err := r.pg.QueryRow(ctx, sql, query.UserID).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func queryApplications(ctx context.Context, q postgres.Querier, sqlQuery string, args ...interface{}) ([]*domain.Application, error) {
	var (
		applications = make([]*domain.Application, 0)

		// scans
		applCreatedAt     *time.Time
		applStatus        *string
		applSpreadsheetID *string
		applLink          *string
		applNo            *int
		digitalAttrs      *any
		financeAttrs      *any
		legalAttrs        *any
	)

	_, err := q.QueryFunc(ctx, sqlQuery, args, []any{
		&applCreatedAt,
		&applStatus,
		&applSpreadsheetID,
		&applLink,
		&applNo,
		&digitalAttrs,
		&financeAttrs,
		&legalAttrs,
	}, func(pgx.QueryFuncRow) error {
		applications = append(applications, &domain.Application{
			CreatedAt:     postgres.Value(applCreatedAt),
			Status:        postgres.Value(applStatus),
			SpreadsheetID: postgres.Value(applSpreadsheetID),
			Link:          postgres.Value(applLink),
			ApplicationNo: postgres.Value(applNo),
			DigitalAttrs:  postgres.Value(digitalAttrs),
			FinanceAttrs:  postgres.Value(financeAttrs),
			LegalAttrs:    postgres.Value(legalAttrs),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return applications, err
}
