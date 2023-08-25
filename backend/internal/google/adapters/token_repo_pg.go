package adapters

import (
	"context"
	"encoding/json"

	"github.com/doodocs/qaztrade/backend/internal/google/domain"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/oauth2"
)

type TokenRepositoryPostgre struct {
	pg *pgxpool.Pool
}

var _ domain.TokenRepository = (*TokenRepositoryPostgre)(nil)

func NewTokenRepositoryPostgre(pg *pgxpool.Pool) *TokenRepositoryPostgre {
	return &TokenRepositoryPostgre{
		pg: pg,
	}
}

func (s *TokenRepositoryPostgre) UpdateToken(ctx context.Context, token *oauth2.Token) error {
	tokenStr, err := encodeToken(token)
	if err != nil {
		return err
	}

	const sql = `
		insert into "oauth2_tokens" ("id", "token") values 
			(1, $1)
		on conflict ("id")
		do update set "token" = excluded."token"
	`
	if _, err := s.pg.Exec(ctx, sql, tokenStr); err != nil {
		return err
	}

	return nil
}

func encodeToken(token *oauth2.Token) (string, error) {
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return "", err
	}

	return string(tokenBytes), nil
}
