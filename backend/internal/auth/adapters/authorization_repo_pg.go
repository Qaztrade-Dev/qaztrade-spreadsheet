package adapters

import (
	"context"
	"encoding/json"

	"github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthorizationRepositoryPostgre struct {
	pg *pgxpool.Pool
}

var _ domain.AuthorizationRepository = (*AuthorizationRepositoryPostgre)(nil)

func NewAuthorizationRepositoryPostgre(pg *pgxpool.Pool) *AuthorizationRepositoryPostgre {
	return &AuthorizationRepositoryPostgre{
		pg: pg,
	}
}

func (r *AuthorizationRepositoryPostgre) SignUp(ctx context.Context, input *domain.SignUpInput) error {
	err := performInTransaction(ctx, r.pg, func(ctx context.Context, tx pgx.Tx) error {
		if err := r.signUp(ctx, tx, input); err != nil {
			return nil
		}

		if err := r.assignRole(ctx, tx, input.UserID, domain.RoleUser); err != nil {
			return nil
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthorizationRepositoryPostgre) signUp(ctx context.Context, tx pgx.Tx, input *domain.SignUpInput) error {
	var (
		hashedPassBytes, _ = bcrypt.GenerateFromPassword([]byte(input.Password), 14)
		hashedPassStr      = string(hashedPassBytes)
	)

	jsonAttrs, err := EncodeUserAttrs(&UserAttrs{OrgName: input.OrgName})
	if err != nil {
		return err
	}

	const sql = `
		insert into "users" 
			(id, email, hashed_password, attrs)
		values
			($1, $2, $3, $4)
	`

	if _, err := tx.Exec(ctx, sql, input.UserID, input.Email, hashedPassStr, jsonAttrs); err != nil {
		return err
	}

	return nil
}

func (r *AuthorizationRepositoryPostgre) assignRole(ctx context.Context, tx pgx.Tx, userID, role string) error {
	const sql = `
		insert into "user_role_bindings"
			(user_id, role_id)
		values
			(
				$1,
				(select id from user_roles where value = $2)
			)
	`

	if _, err := tx.Exec(ctx, sql, userID, role); err != nil {
		return err
	}

	return nil
}

func (r *AuthorizationRepositoryPostgre) GetRoles(ctx context.Context, userID string) ([]string, error) {
	const sql = `
		select
			ur.value
		from user_role_bindings urb
		join user_roles ur on ur.id = urb.role_id
		where
			urb.user_id = $1
	`

	var (
		userRoles []string
	)

	rows, err := r.pg.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var userRole string
		if err := rows.Scan(&userRole); err != nil {
			return nil, err
		}
		userRoles = append(userRoles, userRole)
	}

	return userRoles, nil
}

func (r *AuthorizationRepositoryPostgre) SignIn(ctx context.Context, input *domain.SignInInput) (*domain.User, error) {
	const sql = `
		select 
			u.id,
			u.hashed_password
		from "users" u
		join "user_roles" ur on ur.id = u.role_id
		where 
			u.email = $1
	`

	var (
		userID     string
		hashedPass string
		userRole   string
	)

	err := r.pg.QueryRow(ctx, sql, input.Email).Scan(&userID, &hashedPass, &userRole)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(input.Password)); err != nil {
		return nil, err
	}

	return &domain.User{ID: userID}, nil
}

func (r *AuthorizationRepositoryPostgre) UpdatePassword(ctx context.Context, userID, password string) error {
	var (
		hashedPassBytes, _ = bcrypt.GenerateFromPassword([]byte(password), 14)
		hashedPassStr      = string(hashedPassBytes)
	)

	const sql = `
		update "users" set 
			hashed_password = $2
		where id = $1
	`

	if _, err := r.pg.Exec(ctx, sql, userID, hashedPassStr); err != nil {
		return err
	}

	return nil
}

func (r *AuthorizationRepositoryPostgre) GetOne(ctx context.Context, input *domain.GetQuery) (*domain.User, error) {
	const sql = `
		select 
			u.id,
			ur.value
		from "users" u
		join "user_roles" ur on ur.id = u.role_id
		where 
			u.email = $1
	`

	var (
		userID   string
		userRole string
	)

	err := r.pg.QueryRow(ctx, sql, input.Email).Scan(&userID, &userRole)
	if err != nil {
		return nil, err
	}

	return &domain.User{ID: userID}, nil
}

type UserAttrs struct {
	OrgName string `json:"org_name"`
}

func DecodeUserAttrs(jsonBytes []byte) (*UserAttrs, error) {
	var attrs UserAttrs

	if err := json.Unmarshal(jsonBytes, &attrs); err != nil {
		return nil, err
	}

	return &attrs, nil
}

func EncodeUserAttrs(attrs *UserAttrs) ([]byte, error) {
	return json.Marshal(attrs)
}

type transactionClosure func(ctx context.Context, tx pgx.Tx) error

func performInTransaction(ctx context.Context, pg *pgxpool.Pool, closure transactionClosure) error {
	tx, err := pg.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := closure(ctx, tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
