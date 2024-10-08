package adapters

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/doodocs/qaztrade/backend/internal/manager/domain"
	"github.com/doodocs/qaztrade/backend/pkg/postgres"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mattermost/squirrel"
)

type ManagersRepositoryPostgres struct {
	pg *pgxpool.Pool
}

var _ domain.ManagersRepository = (*ManagersRepositoryPostgres)(nil)

func NewManagersRepositoryPostgres(pg *pgxpool.Pool) *ManagersRepositoryPostgres {
	return &ManagersRepositoryPostgres{
		pg: pg,
	}
}

func (r *ManagersRepositoryPostgres) GetCurrent(ctx context.Context, userID string) (*domain.Manager, error) {
	managers, err := r.getMany(ctx)
	if err != nil {
		return nil, err
	}
	for _, j := range managers {
		j := j
		if j.UserID == userID {
			return j, nil
		}
	}
	return nil, errors.New("current user is not manager")
}

func (r *ManagersRepositoryPostgres) GetMany(ctx context.Context) ([]*domain.Manager, error) {
	managers, err := r.getMany(ctx)
	if err != nil {
		return nil, err
	}

	return managers, nil
}

func getManagersQueryStatement() squirrel.SelectBuilder {
	mainStmt := psql.
		Select(
			"u.id",
			"u.email",
			"u.attrs->>'full_name'",
			"json_agg(ur.value) as roles",
		).
		From("users u").
		Join("user_role_bindings urb ON urb.user_id = u.id").
		Join("user_roles ur ON ur.id = urb.role_id and ur.value <> 'user'").
		GroupBy("u.id").
		OrderBy("u.created_at asc")

	return mainStmt
}

func (r *ManagersRepositoryPostgres) getMany(ctx context.Context) ([]*domain.Manager, error) {
	stmt := getManagersQueryStatement()
	sql, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}

	objects, err := queryManagers(ctx, r.pg, sql, args...)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

func queryManagers(ctx context.Context, q postgres.Querier, sqlQuery string, args ...interface{}) ([]*domain.Manager, error) {
	var (
		objects = make([]*domain.Manager, 0)

		// scans
		tmpUserID   string
		tmpEmail    string
		tmpFullname *string
		tmpRoles    string
	)

	_, err := q.QueryFunc(ctx, sqlQuery, args, []any{
		&tmpUserID,
		&tmpEmail,
		&tmpFullname,
		&tmpRoles,
	}, func(pgx.QueryFuncRow) error {
		var roles []string
		if err := json.Unmarshal([]byte(tmpRoles), &roles); err != nil {
			return err
		}

		objects = append(objects, &domain.Manager{
			UserID:   tmpUserID,
			Email:    tmpEmail,
			Fullname: postgres.Value(tmpFullname),
			Roles:    roles,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	return objects, err
}
