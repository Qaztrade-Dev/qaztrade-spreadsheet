package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mattermost/squirrel"
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

func (r *ApplicationRepositoryPostgre) GetOne(ctx context.Context, query *domain.ApplicationQuery) (*domain.Application, error) {
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

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func (r *ApplicationRepositoryPostgre) getMany(ctx context.Context, query *domain.ApplicationQuery) ([]*domain.Application, error) {
	mainStmt := psql.
		Select(
			"a.created_at",
			"ast.value",
			"a.spreadsheet_id",
			"a.link",
		).
		From("applications a").
		Join("application_statuses ast on ast.id = a.status_id").
		OrderBy("a.created_at desc").
		Limit(query.Limit).Offset(query.Offset)

	if query.ApplicationID != "" {
		mainStmt = mainStmt.Where("a.id = ?", query.ApplicationID)
	}

	sql, args, err := mainStmt.ToSql()
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

func (r *ApplicationRepositoryPostgre) getCount(ctx context.Context, query *domain.ApplicationQuery) (uint64, error) {
	const sql = `
		select 
			count(*)
		from "applications"
	`

	var count uint64
	if err := r.pg.QueryRow(ctx, sql).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

type querier interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
}

func queryApplications(ctx context.Context, q querier, sqlQuery string, args ...interface{}) ([]*domain.Application, error) {
	var (
		applications = make([]*domain.Application, 0)

		// scans
		applCreatedAt     *time.Time
		applStatus        *string
		applSpreadsheetID *string
		applLink          *string
	)

	_, err := q.QueryFunc(ctx, sqlQuery, args, []any{
		&applCreatedAt,
		&applStatus,
		&applSpreadsheetID,
		&applLink,
	}, func(pgx.QueryFuncRow) error {
		applications = append(applications, &domain.Application{
			CreatedAt:     valueFromPointer(applCreatedAt),
			Status:        valueFromPointer(applStatus),
			SpreadsheetID: valueFromPointer(applSpreadsheetID),
			Link:          valueFromPointer(applLink),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return applications, err
}

func valueFromPointer[T any](value *T) T {
	var defaultValue T

	if value == nil {
		return defaultValue
	}
	return *value
}
