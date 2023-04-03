package adapters

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/google/uuid"
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

func (r *AuthorizationRepositoryPostgre) SignUp(ctx context.Context, input *domain.SignUpInput) (*domain.User, error) {
	var (
		userID             = uuid.NewString()
		email              = strings.TrimSpace(input.Email)
		hashedPassBytes, _ = bcrypt.GenerateFromPassword([]byte(input.Password), 14)
		hashedPassStr      = string(hashedPassBytes)
	)

	jsonAttrs, err := EncodeUserAttrs(&UserAttrs{OrgName: input.OrgName})
	if err != nil {
		return nil, err
	}

	const sql = `
		insert into "users" 
			(id, email, hashed_password, attrs, role_id)
		values
			($1, $2, $3, $4,
				(select id from user_roles where value = 'user')
			)
	`

	if _, err := r.pg.Exec(ctx, sql, userID, email, hashedPassStr, jsonAttrs); err != nil {
		return nil, err
	}

	return &domain.User{ID: userID, Role: domain.RoleUser}, nil
}

func (r *AuthorizationRepositoryPostgre) SignIn(ctx context.Context, input *domain.SignInInput) (*domain.User, error) {
	var (
		email = strings.TrimSpace(input.Email)
	)

	const sql = `
		select 
			u.id,
			u.hashed_password,
			ur.value
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

	err := r.pg.QueryRow(ctx, sql, email).Scan(&userID, &hashedPass, &userRole)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(input.Password)); err != nil {
		return nil, err
	}

	return &domain.User{ID: userID, Role: userRole}, nil
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
	var (
		email = strings.TrimSpace(input.Email)
	)

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

	err := r.pg.QueryRow(ctx, sql, email).Scan(&userID, &userRole)
	if err != nil {
		return nil, err
	}

	return &domain.User{ID: userID, Role: userRole}, nil
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
