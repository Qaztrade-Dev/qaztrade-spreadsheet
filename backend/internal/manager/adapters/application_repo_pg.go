package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
	"github.com/doodocs/qaztrade/backend/pkg/postgres"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mattermost/squirrel"
)

type ApplicationRepositoryPostgre struct {
	pg *pgxpool.Pool
}

var _ domain.ApplicationRepository = (*ApplicationRepositoryPostgre)(nil)

func NewApplicationRepositoryPostgres(pg *pgxpool.Pool) *ApplicationRepositoryPostgre {
	return &ApplicationRepositoryPostgre{
		pg: pg,
	}
}

func (r *ApplicationRepositoryPostgre) EditStatus(ctx context.Context, applicationID, statusName string) error {
	const sql = `
		update "applications" set
			status_id = (select id from "application_statuses" where value = $2)
		where 
			id = $1
	`

	if _, err := r.pg.Exec(ctx, sql, applicationID, statusName); err != nil {
		return err
	}
	return nil
}

func (r *ApplicationRepositoryPostgre) GetOne(ctx context.Context, query *domain.GetManyInput) (*domain.Application, error) {
	query.Limit = 1
	query.Offset = 0
	applications, err := r.getMany(ctx, query)
	if err != nil {
		return nil, err
	}

	if len(applications) == 0 {
		return nil, errors.New("not found")
	}

	return applications[0], nil
}

func (r *ApplicationRepositoryPostgre) GetMany(ctx context.Context, query *domain.GetManyInput) (*domain.ApplicationList, error) {
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

func getApplicationQueryStatement(input *domain.GetManyInput) squirrel.SelectBuilder {
	mainStmt := psql.
		Select(
			"a.id",
			"a.no",
			"a.created_at",
			"ast.value",
			"a.spreadsheet_id",
			"a.link",
			"a.sign_document_id",
			"a.sign_at",
			`
			jsonb_build_object(
				'application', a.attrs->'application',
				'sheets_agg', a.attrs->'sheets_agg',
				'sheets', (
					SELECT jsonb_agg(
						jsonb_build_object(
							'title', sub.item->>'title',
							'sheet_id', coalesce(sub.item->>'sheet_id', '0')::double precision,
							'rows', coalesce(sub.item->>'rows', '0')::double precision,
							'expenses', coalesce(sub.item->>'expenses', '0')::double precision,
							'assignee_id_digital', (
								select user_id from assignments where sheet_id = coalesce(sub.item->>'sheet_id', '0')::bigint and type = 'digital'
							),
							'assignee_id_finance', (
								select user_id from assignments where sheet_id = coalesce(sub.item->>'sheet_id', '0')::bigint and type = 'finance'
							),
							'assignee_id_legal', (
								select user_id from assignments where sheet_id = coalesce(sub.item->>'sheet_id', '0')::bigint and type = 'legal'
							)
						)
					)
					FROM jsonb_array_elements(a.attrs->'sheets') sub(item)
				)
			)
			`,
		).
		From("applications a").
		Join("application_statuses ast on ast.id = a.status_id").
		Where("a.no > 0").
		OrderBy("a.no asc")

	if input.BIN != "" {
		mainStmt = mainStmt.Where("a.attrs->'application'->>'bin' = ?", input.BIN)
	}

	if input.CompensationType != "" {
		mainStmt = mainStmt.Where(`EXISTS (
			select 1
			from jsonb_array_elements(a.attrs->'sheets') as j(sheet)
			where sheet->>'title' = ?
		)`, input.CompensationType)
	}

	if !input.SignedAtFrom.IsZero() {
		timeFromStr := input.SignedAtFrom.
			Truncate(time.Second).
			Truncate(time.Minute).
			Truncate(time.Hour * 24).Format(time.DateOnly)
		mainStmt = mainStmt.Where("date(a.sign_at) >= ?", timeFromStr)
	}

	if !input.SignedAtUntil.IsZero() {
		timeUntilStr := input.SignedAtUntil.
			Truncate(time.Second).
			Truncate(time.Minute).
			Truncate(time.Hour * 24).Format(time.DateOnly)
		mainStmt = mainStmt.Where("date(a.sign_at) <= ?", timeUntilStr)
	}

	if input.ApplicationID != "" {
		mainStmt = mainStmt.Where("a.id = ?", input.ApplicationID)
	}

	if input.CompanyName != "" {
		mainStmt = mainStmt.Where("a.attrs->'application'->>'from' ilike ?", "%"+input.CompanyName+"%")
	}

	if input.ApplicationNo != 0 {
		mainStmt = mainStmt.Where("a.no = ?", input.ApplicationNo)
	}

	return mainStmt
}

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func (r *ApplicationRepositoryPostgre) getMany(ctx context.Context, input *domain.GetManyInput) ([]*domain.Application, error) {
	stmt := getApplicationQueryStatement(input).Limit(input.Limit).Offset(input.Offset)
	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}

	applications, err := queryApplications(ctx, r.pg, sql, args...)
	if err != nil {
		return nil, err
	}

	for i := range applications {
		applications[i].Link = fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/edit", applications[i].SpreadsheetID)
	}

	return applications, nil
}

func (r *ApplicationRepositoryPostgre) getCount(ctx context.Context, query *domain.GetManyInput) (uint64, error) {
	stmt := getApplicationQueryStatement(query)
	sql, args, err := psql.Select("count(*)").FromSelect(stmt, "q").ToSql()
	if err != nil {
		return 0, fmt.Errorf("countRows %w", err)
	}

	var tmp uint64
	err = r.pg.QueryRow(ctx, sql, args...).Scan(&tmp)
	if err != nil {
		err = fmt.Errorf("count query %w", err)
	}

	return tmp, err
}

func queryApplications(ctx context.Context, q postgres.Querier, sqlQuery string, args ...interface{}) ([]*domain.Application, error) {
	var (
		applications = make([]*domain.Application, 0)

		// scans
		applID             *string
		applNo             *int
		applCreatedAt      *time.Time
		applStatus         *string
		applSpreadsheetID  *string
		applLink           *string
		applSignDocumentID *string
		applSignAt         *time.Time
		applAttrs          *interface{}
	)

	_, err := q.QueryFunc(ctx, sqlQuery, args, []any{
		&applID,
		&applNo,
		&applCreatedAt,
		&applStatus,
		&applSpreadsheetID,
		&applLink,
		&applSignDocumentID,
		&applSignAt,
		&applAttrs,
	}, func(pgx.QueryFuncRow) error {
		applications = append(applications, &domain.Application{
			ID:             postgres.Value(applID),
			No:             postgres.Value(applNo),
			CreatedAt:      postgres.Value(applCreatedAt),
			Status:         postgres.Value(applStatus),
			SpreadsheetID:  postgres.Value(applSpreadsheetID),
			Link:           postgres.Value(applLink),
			SignDocumentID: postgres.Value(applSignDocumentID),
			SignedAt:       postgres.Value(applSignAt),
			Attrs:          postgres.Value(applAttrs),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return applications, err
}
