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

	spreadsheetSvc, err := adapters.NewSpreadsheetService(ctx, deps.credentials, deps.adminAccount)
	if err != nil {
		panic(err)
	}

	var (
		applicationRepo = adapters.NewApplicationRepositoryPostgres(deps.pg)
		managersRepo    = adapters.NewManagersRepositoryPostgres(deps.pg)
		signSvc         = adapters.NewSigningServiceDoodocs(deps.signUrlBase, deps.signLogin, deps.signPassword)
	)

	svc := service.NewService(spreadsheetSvc, applicationRepo, signSvc, managersRepo)
	return svc
}

type Option func(*dependencies)

type dependencies struct {
	credentials []byte
	pg          *pgxpool.Pool

	signUrlBase  string
	signLogin    string
	signPassword string

	adminAccount string
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

func WithSignCredentials(signUrlBase, signLogin, signPassword string) Option {
	return func(d *dependencies) {
		d.signUrlBase = signUrlBase
		d.signLogin = signLogin
		d.signPassword = signPassword
	}
}

func WithAdmin(input string) Option {
	return func(d *dependencies) {
		d.adminAccount = input
	}
}
