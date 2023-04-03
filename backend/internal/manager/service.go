package manager

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/manager/adapters"
	"github.com/doodocs/qaztrade/backend/internal/manager/service"
	"github.com/jackc/pgx/v4/pgxpool"
)

func MakeService(ctx context.Context, opts ...Option) service.Service {
	deps := &dependencies{}
	deps.setDefaults()
	for _, opt := range opts {
		opt(deps)
	}

	spreadsheetSvc, err := adapters.NewSpreadsheetService(ctx, deps.credentials)
	if err != nil {
		panic(err)
	}

	applicationRepo := adapters.NewApplicationRepositoryPostgre(deps.pg)

	svc := service.NewService(spreadsheetSvc, applicationRepo)
	return svc
}

type Option func(*dependencies)

type dependencies struct {
	credentials []byte
	pg          *pgxpool.Pool
}

func (d *dependencies) setDefaults() {
	// pass
}

func WithPostgre(pg *pgxpool.Pool) Option {
	return func(d *dependencies) {
		d.pg = pg
	}
}

func WithCredentials(credentials []byte) Option {
	return func(d *dependencies) {
		d.credentials = credentials
	}
}
