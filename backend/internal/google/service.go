package google

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/google/adapters"
	"github.com/doodocs/qaztrade/backend/internal/google/service"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/oauth2"
)

func MakeService(ctx context.Context, opts ...Option) service.Service {
	deps := &dependencies{}
	deps.setDefaults()
	for _, opt := range opts {
		opt(deps)
	}

	var (
		tokenRepo = adapters.NewTokenRepositoryPostgre(deps.pg)
		svc       = service.NewService(deps.config, tokenRepo)
	)

	return svc
}

type Option func(*dependencies)

type dependencies struct {
	pg     *pgxpool.Pool
	config *oauth2.Config
}

func (d *dependencies) setDefaults() {
	// pass
}

func WithPostgre(pg *pgxpool.Pool) Option {
	return func(d *dependencies) {
		d.pg = pg
	}
}

func WithOAuthConfig(config *oauth2.Config) Option {
	return func(d *dependencies) {
		d.config = config
	}
}
