package adapters

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/domain"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepositoryPostgre struct {
	pg *pgxpool.Pool
}

var _ domain.UserRepository = (*UserRepositoryPostgre)(nil)

func NewUserRepositoryPostgre(pg *pgxpool.Pool) *UserRepositoryPostgre {
	return &UserRepositoryPostgre{
		pg: pg,
	}
}

func (r *UserRepositoryPostgre) Get(ctx context.Context, userID string) (*domain.User, error) {
	const sql = `
		select 
			attrs
		from "users"
		where 
			id = $1
	`

	var attrs UserAttrs
	err := r.pg.QueryRow(ctx, sql, userID).Scan(&attrs)
	if err != nil {
		return nil, err
	}

	return DecodeUserAttrs(userID, &attrs), nil
}

type UserAttrs struct {
	OrgName string `json:"org_name"`
}

func DecodeUserAttrs(userID string, attrs *UserAttrs) *domain.User {
	return &domain.User{
		ID:      userID,
		OrgName: attrs.OrgName,
	}
}
