package auth

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/auth/adapters"
	"github.com/doodocs/qaztrade/backend/internal/auth/service"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
)

func MakeService(ctx context.Context, opts ...Option) service.Service {
	deps := &dependencies{}
	deps.setDefaults()
	for _, opt := range opts {
		opt(deps)
	}

	var (
		authRepo  = adapters.NewAuthorizationRepositoryPostgre(deps.pg)
		credsRepo = adapters.NewCredentialsRepositoryJWT(deps.jwtcli)
		emailSvc  = adapters.NewEmailServiceGmail(deps.mailEmail, deps.mailPassword)
	)

	svc := service.NewService(authRepo, credsRepo, emailSvc)
	return svc
}

type Option func(*dependencies)

type dependencies struct {
	pg     *pgxpool.Pool
	jwtcli *jwt.Client

	mailEmail, mailPassword string
}

func (d *dependencies) setDefaults() {
	// pass
}

func WithPostgre(pg *pgxpool.Pool) Option {
	return func(d *dependencies) {
		d.pg = pg
	}
}

func WithJWT(jwtcli *jwt.Client) Option {
	return func(d *dependencies) {
		d.jwtcli = jwtcli
	}
}

func WithMail(mailEmail, mailPassword string) Option {
	return func(d *dependencies) {
		d.mailEmail = mailEmail
		d.mailPassword = mailPassword
	}
}
