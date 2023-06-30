package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mattermost/squirrel"
)

type AssignmentsRepositoryPostgres struct {
	pg *pgxpool.Pool
}

var _ domain.AssignmentsRepository = (*AssignmentsRepositoryPostgres)(nil)

func NewAssignmentsRepositoryPostgres(pg *pgxpool.Pool) *AssignmentsRepositoryPostgres {
	return &AssignmentsRepositoryPostgres{
		pg: pg,
	}
}

func (r *AssignmentsRepositoryPostgres) GetInfo(ctx context.Context, input *domain.GetInfoInput) (*domain.AssignmentsInfo, error) {
	stmt := getAssignmentsInfoQueryStatement(input)
	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}

	var total, completed uint64
	err = r.pg.QueryRow(ctx, sql, args...).Scan(&total, &completed)
	if err != nil {
		return nil, err
	}

	output := &domain.AssignmentsInfo{
		Total:     total,
		Completed: completed,
	}

	return output, nil
}

func getAssignmentsInfoQueryStatement(input *domain.GetInfoInput) squirrel.SelectBuilder {
	mainStmt := psql.
		Select(
			"count(*) as total",
			"count(*) filter (where ass.is_completed) as completed",
		).
		From("assignments ass").
		Join("users u on u.id = ass.user_id")

	if input.UserID != nil {
		mainStmt = mainStmt.Where("u.id = ?", *input.UserID)
	}

	return mainStmt
}

func (r *AssignmentsRepositoryPostgres) GetMany(ctx context.Context, input *domain.GetManyInput) (*domain.AssignmentsList, error) {
	total, err := r.getCount(ctx, input)
	if err != nil {
		return nil, err
	}

	objects, err := r.getMany(ctx, input)
	if err != nil {
		return nil, err
	}

	result := &domain.AssignmentsList{
		Total:   int(total),
		Objects: objects,
	}

	return result, nil
}

func getAssignmentsQueryStatement(input *domain.GetManyInput) squirrel.SelectBuilder {
	mainStmt := psql.
		Select(
			"ass.id",
			"app.attrs->'application'->>'from'",
			"app.attrs->'application'->>'bin'",
			"ass.sheet_title",
			"ass.sheet_id",
			"ass.type",
			"app.link",
			"u.attrs->>'name'",
			"ass.rows_from",
			"ass.rows_until",
			"coalesce(assres.total_completed, 0)",
			"ass.is_completed",
			"ass.completed_at",
		).
		From("assignments ass").
		Join("applications app on app.id = ass.application_id").
		Join("users u on u.id = ass.user_id").
		LeftJoin("assignment_results assres on assres.id = ass.last_result_id").
		OrderBy("ass.created_at asc")

	if input.UserID != nil {
		mainStmt = mainStmt.Where("u.id = ?", *input.UserID)
	}

	return mainStmt
}

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func (r *AssignmentsRepositoryPostgres) getMany(ctx context.Context, input *domain.GetManyInput) ([]*domain.AssignmentView, error) {
	stmt := getAssignmentsQueryStatement(input).Limit(input.Limit).Offset(input.Offset)
	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}

	objects, err := queryAssignmentViews(ctx, r.pg, sql, args...)
	if err != nil {
		return nil, err
	}

	for i := range objects {
		objects[i].Link = fmt.Sprintf("%s#gid=%v", objects[i].Link, objects[i].SheetID)
		objects[i].RowsTotal = (objects[i].RowsUntil - objects[i].RowsFrom) + 1
		objects[i].RowsFrom += 3
		objects[i].RowsUntil += 3
	}

	return objects, nil
}

func (r *AssignmentsRepositoryPostgres) getCount(ctx context.Context, query *domain.GetManyInput) (uint64, error) {
	stmt := getAssignmentsQueryStatement(query)
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

type querier interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	QueryFunc(ctx context.Context, sql string, args []interface{}, scans []interface{}, f func(pgx.QueryFuncRow) error) (pgconn.CommandTag, error)
}

func queryAssignmentViews(ctx context.Context, q querier, sqlQuery string, args ...interface{}) ([]*domain.AssignmentView, error) {
	var (
		objects = make([]*domain.AssignmentView, 0)

		// scans
		tmpID             *int
		tmpApplicantName  *string
		tmpApplicantBIN   *string
		tmpSheetTitle     *string
		tmpSheetID        *uint64
		tmpAssignmentType *string
		tmpLink           *string
		tmpAssigneeName   *string
		tmpRowsFrom       *int
		tmpRowsUntil      *int
		tmpRowsCompleted  *int
		tmpIsCompleted    *bool
		tmpCompletedAt    *time.Time
	)

	_, err := q.QueryFunc(ctx, sqlQuery, args, []any{
		&tmpID,
		&tmpApplicantName,
		&tmpApplicantBIN,
		&tmpSheetTitle,
		&tmpSheetID,
		&tmpAssignmentType,
		&tmpLink,
		&tmpAssigneeName,
		&tmpRowsFrom,
		&tmpRowsUntil,
		&tmpRowsCompleted,
		&tmpIsCompleted,
		&tmpCompletedAt,
	}, func(pgx.QueryFuncRow) error {
		objects = append(objects, &domain.AssignmentView{
			ID:             valueFromPointer(tmpID),
			ApplicantName:  valueFromPointer(tmpApplicantName),
			ApplicantBIN:   valueFromPointer(tmpApplicantBIN),
			SheetTitle:     valueFromPointer(tmpSheetTitle),
			SheetID:        valueFromPointer(tmpSheetID),
			AssignmentType: valueFromPointer(tmpAssignmentType),
			Link:           valueFromPointer(tmpLink),
			AssigneeName:   valueFromPointer(tmpAssigneeName),
			RowsFrom:       valueFromPointer(tmpRowsFrom),
			RowsUntil:      valueFromPointer(tmpRowsUntil),
			RowsCompleted:  valueFromPointer(tmpRowsCompleted),
			IsCompleted:    valueFromPointer(tmpIsCompleted),
			CompletedAt:    valueFromPointer(tmpCompletedAt),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return objects, err
}

func valueFromPointer[T any](value *T) T {
	var defaultValue T

	if value == nil {
		return defaultValue
	}
	return *value
}
