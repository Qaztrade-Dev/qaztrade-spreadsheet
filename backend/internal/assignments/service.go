package assignments

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/assignments/adapters"
	"github.com/doodocs/qaztrade/backend/internal/assignments/service"
	"github.com/jackc/pgx/v4/pgxpool"
)

func MakeService(ctx context.Context, opts ...Option) service.Service {
	deps := &dependencies{}
	deps.setDefaults()
	for _, opt := range opts {
		opt(deps)
	}

	var (
		assignmentsRepo = adapters.NewAssignmentsRepositoryPostgres(deps.pg)
	)

	svc := service.NewService(assignmentsRepo)
	return svc
}

type Option func(*dependencies)

type dependencies struct {
	pg *pgxpool.Pool
}

func (d *dependencies) setDefaults() {
	// pass
}

func WithPostgres(pg *pgxpool.Pool) Option {
	return func(d *dependencies) {
		d.pg = pg
	}
}
