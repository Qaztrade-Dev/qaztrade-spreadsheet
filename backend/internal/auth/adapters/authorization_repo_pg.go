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

func (r *AuthorizationRepositoryPostgre) SignUp(ctx context.Context, input *domain.SignUpInput) (string, error) {
	var (
		userID             = uuid.NewString()
		email              = strings.TrimSpace(input.Email)
		hashedPassBytes, _ = bcrypt.GenerateFromPassword([]byte(input.Password), 14)
		hashedPassStr      = string(hashedPassBytes)
	)

	jsonAttrs, err := EncodeUserAttrs(&UserAttrs{OrgName: input.OrgName})
	if err != nil {
		return "", err
	}

	const sql = `
		insert into "users" 
			(id, email, hashed_password, attrs)
		values
			($1, $2, $3, $4)
	`

	if _, err := r.pg.Exec(ctx, sql, userID, email, hashedPassStr, jsonAttrs); err != nil {
		return "", err
	}

	return userID, nil
}

func (r *AuthorizationRepositoryPostgre) SignIn(ctx context.Context, input *domain.SignInInput) (string, error) {
	var (
		email              = strings.TrimSpace(input.Email)
		hashedPassBytes, _ = bcrypt.GenerateFromPassword([]byte(input.Password), 14)
		hashedPassStr      = string(hashedPassBytes)
	)

	const sql = `
		select 
			id
		from "users"
		where 
			email = $1 and hashed_password = $2
	`

	var userID string
	err := r.pg.QueryRow(ctx, sql, email, hashedPassStr).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
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

func (r *AuthorizationRepositoryPostgre) GetOne(ctx context.Context, input *domain.GetQuery) (string, error) {
	var (
		email = strings.TrimSpace(input.Email)
	)

	const sql = `
		select 
			id
		from "users"
		where 
			email = $1
	`

	var userID string
	err := r.pg.QueryRow(ctx, sql, email).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
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
