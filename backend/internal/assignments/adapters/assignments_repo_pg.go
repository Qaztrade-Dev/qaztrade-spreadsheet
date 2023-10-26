package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/doodocs/qaztrade/backend/internal/assignments/domain"
	"github.com/doodocs/qaztrade/backend/pkg/postgres"
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

func (r *AssignmentsRepositoryPostgres) LockApplications(ctx context.Context) (int, error) {
	var batchID int

	err := postgres.InTransaction(ctx, r.pg, func(ctx context.Context, tx pgx.Tx) error {
		tmpBatchID, err := r.createBatch(ctx, tx)
		if err != nil {
			return err
		}

		if err := r.assignBatchApplications(ctx, tx, tmpBatchID); err != nil {
			return err
		}

		batchID = tmpBatchID
		return nil
	})

	if err != nil {
		return batchID, err
	}
	return batchID, nil
}

func (r *AssignmentsRepositoryPostgres) createBatch(ctx context.Context, q postgres.Querier) (int, error) {
	const sql = `
		insert into "batches" default values returning id;
	`

	var batchID int
	if err := q.QueryRow(ctx, sql).Scan(&batchID); err != nil {
		return 0, err
	}

	return batchID, nil
}

func (r *AssignmentsRepositoryPostgres) assignBatchApplications(ctx context.Context, q postgres.Querier, batchID int) error {
	const sql = `
		insert into "batch_applications" 
			(batch_id, application_id)
		select
			$1, id
		from applications
		where is_signed
	`

	if _, err := q.Exec(ctx, sql, batchID); err != nil {
		return err
	}

	return nil
}

func (r *AssignmentsRepositoryPostgres) CreateAssignments(ctx context.Context, inputs []*domain.AssignmentInput) error {
	err := postgres.InTransaction(ctx, r.pg, func(ctx context.Context, tx pgx.Tx) error {
		for _, input := range inputs {
			if err := r.createAssignment(ctx, tx, input); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (r *AssignmentsRepositoryPostgres) createAssignment(ctx context.Context, q postgres.Querier, input *domain.AssignmentInput) error {
	const sql = `
		insert into "assignments" 
			(
				user_id,
				application_id,
				type,
				sheet_title,
				sheet_id,
				total_rows,
				total_sum
			)
		values
			($1, $2, $3, $4, $5, $6, $7)
	`

	if _, err := q.Exec(ctx, sql,
		input.ManagerID,
		input.ApplicationID,
		input.AssignmentType,
		input.SheetTitle,
		input.SheetID,
		input.TotalRows,
		input.TotalSum,
	); err != nil {
		return err
	}

	return nil
}

func (r *AssignmentsRepositoryPostgres) GetManagerIDs(ctx context.Context, role string) ([]string, error) {
	const sql = `
		select
			u.id
		from users u
		join user_role_bindings urb on urb.user_id = u.id
		join user_roles ur on ur.id = urb.role_id
		where 
			ur.value = $1		
	`

	objects, err := queryStrings(ctx, r.pg, sql, role)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

func queryStrings(ctx context.Context, q postgres.Querier, sqlQuery string, args ...interface{}) ([]string, error) {
	var (
		objects = make([]string, 0)

		// scans
		tmpStr *string
	)

	_, err := q.QueryFunc(ctx, sqlQuery, args,
		[]any{&tmpStr},
		func(pgx.QueryFuncRow) error {
			objects = append(objects, postgres.Value(tmpStr))
			return nil
		})
	if err != nil {
		return nil, err
	}

	return objects, err
}

func (r *AssignmentsRepositoryPostgres) GetSheets(ctx context.Context, batchID int, sheetTable string) ([]*domain.Sheet, error) {
	stmt := getSheetsQueryStatement(batchID, sheetTable)
	sql, args, err := stmt.ToSql()

	if err != nil {
		return nil, err
	}

	objects, err := querySheets(ctx, r.pg, sql, args...)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

func getSheetsQueryStatement(batchID int, sheetTable string) squirrel.SelectBuilder {
	mainStmt := psql.
		Select(
			"a.id",
			"e.sheet_title",
			"safe_cast_to_int(s.value->>'sheet_id')",
			"e.total_rows",
			"e.total_sum",
		).
		From("applications a").
		CrossJoin("jsonb_array_elements(a.attrs -> 'sheets') as s").
		Join("expenses_agg e on e.id = a.id and e.sheet_title = s.value ->> 'title'")
		// Join(sheetTable + " e on e.spreadsheet_id = a.spreadsheet_id and e.sheet_title = s.value ->> 'title'").
		// Join("applicants_info_view info on info.id = a.id")
		// // Where("a.id in (select application_id from batch_applications where batch_id = ?)", batchID)

	return mainStmt
}

func querySheets(ctx context.Context, q postgres.Querier, sqlQuery string, args ...interface{}) ([]*domain.Sheet, error) {
	var (
		objects = make([]*domain.Sheet, 0)

		// scans
		tmpApplicationID *string
		tmpSheetTitle    *string
		tmpSheetID       *uint64
		tmpTotalRows     *uint64
		tmpTotalSum      *float64
	)

	_, err := q.QueryFunc(ctx, sqlQuery, args, []any{
		&tmpApplicationID,
		&tmpSheetTitle,
		&tmpSheetID,
		&tmpTotalRows,
		&tmpTotalSum,
	}, func(pgx.QueryFuncRow) error {
		objects = append(objects, &domain.Sheet{
			ApplicationID: postgres.Value(tmpApplicationID),
			SheetTitle:    postgres.Value(tmpSheetTitle),
			SheetID:       postgres.Value(tmpSheetID),
			TotalRows:     postgres.Value(tmpTotalRows),
			TotalSum:      postgres.Value(tmpTotalSum),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return objects, err
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

func (r *AssignmentsRepositoryPostgres) InsertAssignmentResult(ctx context.Context, assignmentID uint64, total uint64) error {
	if err := postgres.InTransaction(ctx, r.pg, func(ctx context.Context, tx pgx.Tx) error {
		var query1 = `insert into "assignment_results" (assignment_id, total_completed) values ($1, $2) returning id`
		var assResID string
		if err := tx.QueryRow(ctx, query1, assignmentID, total).Scan(&assResID); err != nil {
			return fmt.Errorf("insert: %w", err)
		}

		var query2 = `update "assignments" set last_result_id = $2 where id = $1`
		if _, err := tx.Exec(ctx, query2, assignmentID, assResID); err != nil {
			return fmt.Errorf("update last_result_id: %w", err)
		}

		assignment, err := getOne(ctx, tx, &domain.GetManyInput{
			AssignmentID: &assignmentID,
		})
		if err != nil {
			return fmt.Errorf("getOne: %w", err)
		}

		if assignment.TotalRows == assignment.RowsCompleted {
			var query3 = `update "assignments" set is_completed = true, completed_at = now() where id = $1`
			if _, err := tx.Exec(ctx, query3, assignmentID); err != nil {
				return fmt.Errorf("update is_completed: %w", err)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *AssignmentsRepositoryPostgres) ChangeAssignee(ctx context.Context, input *domain.ChangeAssigneeInput) error {
	const sql = `
		update "assignments"
		set
			user_id=$2
		where
			id=$1
	`
	if _, err := r.pg.Exec(ctx, sql, input.AssignmentID, input.UserID); err != nil {
		return err
	}

	return nil
}

func (r *AssignmentsRepositoryPostgres) GetMany(ctx context.Context, input *domain.GetManyInput) (*domain.AssignmentsList, error) {
	total, err := r.getCount(ctx, input)
	if err != nil {
		return nil, err
	}

	objects, err := getMany(ctx, r.pg, input)
	if err != nil {
		return nil, err
	}

	result := &domain.AssignmentsList{
		Total:   int(total),
		Objects: objects,
	}

	return result, nil
}

func (r *AssignmentsRepositoryPostgres) GetOne(ctx context.Context, input *domain.GetManyInput) (*domain.AssignmentView, error) {
	return getOne(ctx, r.pg, input)
}

func getOne(ctx context.Context, querier postgres.Querier, input *domain.GetManyInput) (*domain.AssignmentView, error) {
	input.Limit = 1
	objects, err := getMany(ctx, querier, input)
	if err != nil {
		return nil, err
	}

	if len(objects) == 0 {
		return nil, domain.ErrAssignmentNotFound
	}

	return objects[0], nil
}

func getAssignmentsQueryStatement(input *domain.GetManyInput) squirrel.SelectBuilder {
	const inass = `(
		select 
			coalesce(inappst.value, 'manager_reviewing') as status,
			inass.application_id,
			inass.sheet_title,
			inass.type
		from assignments inass
		left join application_statuses inappst on inappst.id = inass.resolution_status_id
	)`

	mainStmt := psql.
		Select(
			"app.id",
			"ass.id",
			"app.no",
			"app.attrs->'application'->>'from'",
			"app.attrs->'application'->>'bin'",
			"app.spreadsheet_id",
			"ass.sheet_title",
			"ass.sheet_id",
			"ass.type",
			"app.link",
			"app.sign_link",
			"u.attrs->>'full_name'",
			"u.id",
			"ass.total_rows",
			"ass.total_sum",
			"coalesce(assres.total_completed, 0)",
			"ass.is_completed",
			"ass.completed_at",
			"coalesce(appst.value, 'manager_reviewing')",
			"ass.resolved_at",
			"ass.countdown_duration",
			"digital.status",
			"finance.status",
			"legal.status",
		).
		From("assignments ass").
		Join("applications app on app.id = ass.application_id").
		LeftJoin("application_statuses appst on appst.id = ass.resolution_status_id").
		Join("users u on u.id = ass.user_id").
		LeftJoin("assignment_results assres on assres.id = ass.last_result_id").
		LeftJoin(inass+" digital on digital.application_id = ass.application_id and digital.sheet_title = ass.sheet_title and digital.type = 'digital'").
		LeftJoin(inass+" finance on finance.application_id = ass.application_id and finance.sheet_title = ass.sheet_title and finance.type = 'finance'").
		LeftJoin(inass+" legal on legal.application_id = ass.application_id and legal.sheet_title = ass.sheet_title and legal.type = 'legal'").
		OrderBy("app.no asc", "ass.type asc")

	if input.AssigneeID != nil {
		mainStmt = mainStmt.Where("u.id = ?", *input.AssigneeID)
	}

	if input.AssignmentID != nil {
		mainStmt = mainStmt.Where("ass.id = ?", *input.AssignmentID)
	}

	if input.IsCompleted != nil {
		mainStmt = mainStmt.Where("ass.is_completed = ?", *input.IsCompleted)
	}

	if input.CompanyName != nil {
		mainStmt = mainStmt.Where("app.attrs->'application'->>'from' ilike ?", "%"+(*input.CompanyName)+"%")
	}

	if input.ApplicationNo != nil {
		mainStmt = mainStmt.Where("app.no = ?", *input.ApplicationNo)
	}

	if input.AssignmentType != nil {
		mainStmt = mainStmt.Where("ass.type = ?", *input.AssignmentType)
	}

	return mainStmt
}

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func getMany(ctx context.Context, querier postgres.Querier, input *domain.GetManyInput) ([]*domain.AssignmentView, error) {
	stmt := getAssignmentsQueryStatement(input)

	if input.Limit != 0 {
		stmt = stmt.Limit(input.Limit)
	}

	if input.Offset != 0 {
		stmt = stmt.Offset(input.Offset)
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}

	objects, err := queryAssignmentViews(ctx, querier, sql, args...)
	if err != nil {
		return nil, err
	}

	for i := range objects {
		objects[i].Link = fmt.Sprintf("%s#gid=%v", objects[i].Link, objects[i].SheetID)
		objects[i].SignLink = fmt.Sprintf("https://link.doodocs.kz/%s", objects[i].SignLink)
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

func queryAssignmentViews(ctx context.Context, q postgres.Querier, sqlQuery string, args ...interface{}) ([]*domain.AssignmentView, error) {
	var (
		objects = make([]*domain.AssignmentView, 0)

		// scans
		tmpApplicationID     *string
		tmpAssignmentID      *uint64
		tmpID                *int64
		tmpApplicantName     *string
		tmpApplicantBIN      *string
		tmpSpreadsheetID     *string
		tmpSheetTitle        *string
		tmpSheetID           *uint64
		tmpAssignmentType    *string
		tmpLink              *string
		tmpSignLink          *string
		tmpAssigneeName      *string
		tmpAssigneeID        *string
		tmpTotalRows         *int
		tmpTotalSum          *int
		tmpRowsCompleted     *int
		tmpIsCompleted       *bool
		tmpCompletedAt       *time.Time
		tmpResolutionStatus  *string
		tmpResolvedAt        *time.Time
		tmpCountdownDuration *time.Duration
		tmpDigitalStatus     *string
		tmpFinanaceStatus    *string
		tmpLegalStatus       *string
	)

	_, err := q.QueryFunc(ctx, sqlQuery, args, []any{
		&tmpApplicationID,
		&tmpAssignmentID,
		&tmpID,
		&tmpApplicantName,
		&tmpApplicantBIN,
		&tmpSpreadsheetID,
		&tmpSheetTitle,
		&tmpSheetID,
		&tmpAssignmentType,
		&tmpLink,
		&tmpSignLink,
		&tmpAssigneeName,
		&tmpAssigneeID,
		&tmpTotalRows,
		&tmpTotalSum,
		&tmpRowsCompleted,
		&tmpIsCompleted,
		&tmpCompletedAt,
		&tmpResolutionStatus,
		&tmpResolvedAt,
		&tmpCountdownDuration,
		&tmpDigitalStatus,
		&tmpFinanaceStatus,
		&tmpLegalStatus,
	}, func(pgx.QueryFuncRow) error {
		objects = append(objects, &domain.AssignmentView{
			ApplicationID:     postgres.Value(tmpApplicationID),
			AssignmentID:      postgres.Value(tmpAssignmentID),
			ID:                postgres.Value(tmpID),
			ApplicantName:     postgres.Value(tmpApplicantName),
			ApplicantBIN:      postgres.Value(tmpApplicantBIN),
			SpreadsheetID:     postgres.Value(tmpSpreadsheetID),
			SheetTitle:        postgres.Value(tmpSheetTitle),
			SheetID:           postgres.Value(tmpSheetID),
			AssignmentType:    postgres.Value(tmpAssignmentType),
			Link:              postgres.Value(tmpLink),
			SignLink:          postgres.Value(tmpSignLink),
			AssigneeName:      postgres.Value(tmpAssigneeName),
			AssigneeID:        postgres.Value(tmpAssigneeID),
			TotalRows:         postgres.Value(tmpTotalRows),
			TotalSum:          postgres.Value(tmpTotalSum),
			RowsCompleted:     postgres.Value(tmpRowsCompleted),
			IsCompleted:       postgres.Value(tmpIsCompleted),
			CompletedAt:       postgres.Value(tmpCompletedAt),
			ResolutionStatus:  postgres.Value(tmpResolutionStatus),
			ResolvedAt:        postgres.Value(tmpResolvedAt),
			CountdownDuration: postgres.Value(tmpCountdownDuration),
			DigitalStatus:     postgres.Value(tmpDigitalStatus),
			FinanceStatus:     postgres.Value(tmpFinanaceStatus),
			LegalStatus:       postgres.Value(tmpLegalStatus),
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return objects, err
}

func (r *AssignmentsRepositoryPostgres) UpdateAssignees(ctx context.Context, inputs []*domain.AssignmentInput) error {
	if err := postgres.InTransaction(ctx, r.pg, func(ctx context.Context, tx pgx.Tx) error {
		var query = `update "assignments" set user_id = $1 where id = $2`

		for _, input := range inputs {
			_, err := tx.Exec(ctx, query, input.ManagerID, int(input.AssignmentID))
			if err != nil {
				return fmt.Errorf("update user_id: %w", err)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *AssignmentsRepositoryPostgres) SetResolution(ctx context.Context, input *domain.SetResolutionInput) error {
	stmt := psql.Update("assignments").Where("id = ?", input.AssignmentID)
	stmt = stmt.Set("resolved_at", input.ResolvedAt)
	stmt = stmt.Set("countdown_duration", input.CountdownDuration)

	if input.ResolutionStatus != "" {
		stmt = stmt.Set("resolution_status_id", squirrel.Expr("(select id from application_statuses where value = ?)", input.ResolutionStatus))
	}

	sql, args, err := stmt.ToSql()
	if err != nil {
		return err
	}

	if _, err := r.pg.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}

func (r *AssignmentsRepositoryPostgres) AllAssignmentsStatusEq(ctx context.Context, applicationID, statusName string) (bool, error) {
	const sql = `
		select
			$2 = all(
				select
					ast.value
				from assignments ass
				join application_statuses ast on ast.id = ass.resolution_status_id
				where 
					ass.application_id = $1 and
					ass.resolution_status_id is not null
			)
	`

	var all bool
	if err := r.pg.QueryRow(ctx, sql, applicationID, statusName).Scan(&all); err != nil {
		return false, err
	}

	return all, nil
}

func (r *AssignmentsRepositoryPostgres) UpdateStatus(ctx context.Context, input *domain.UpdateStatusInput) error {
	const sql = `
		update "assignments" set
			resolution_status_id = (select id from application_statuses where value = $2)
		where 
			id = $1
	`

	if _, err := r.pg.Exec(ctx, sql, input.AssignmentID, input.StatusName); err != nil {
		return err
	}

	return nil
}
