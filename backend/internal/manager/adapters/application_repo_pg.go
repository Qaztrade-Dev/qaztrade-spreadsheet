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

func (r *ApplicationRepositoryPostgre) IsManagerAssigned(ctx context.Context, applicationID, userID string) (bool, error) {
	const sql = `
		select exists (
			select 1
			from "assignments"
			where 
				application_id = $1 and
				user_id = $2
		)
	`

	var tmp bool
	if err := r.pg.QueryRow(ctx, sql, applicationID, userID).Scan(&tmp); err != nil {
		return false, err
	}

	return tmp, nil
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
	const withQuery = `
	with "application_progress" as (
		SELECT distinct on (ass.application_id)
			ass.application_id,
			ass.total_rows,
			ass.total_sum,
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
	)`

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
			"a.attrs->'application'",
			"apr.digital",
			"apr.finance",
			"apr.legal",
		).Prefix(withQuery).
		From("applications a").
		Join("application_statuses ast on ast.id = a.status_id").
		Join("application_progress apr on apr.application_id = a.id").
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
	stmt := getApplicationQueryStatement(input)

	if input.Limit > 0 {
		stmt = stmt.Limit(input.Limit)
	}

	if input.Offset > 0 {
		stmt = stmt.Offset(input.Offset)
	}

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
		attrsDigital       *interface{}
		attrsFinance       *interface{}
		attrsLegal         *interface{}
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
		&attrsDigital,
		&attrsFinance,
		&attrsLegal,
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
			AttrsDigital:   postgres.Value(attrsDigital),
			AttrsFinance:   postgres.Value(attrsFinance),
			AttrsLegal:     postgres.Value(attrsLegal),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return applications, err
}
