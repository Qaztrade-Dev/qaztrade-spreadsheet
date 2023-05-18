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

	spreadsheetStorage, err := adapters.NewStorageS3(ctx, deps.s3AccessKey, deps.s3SecretKey, deps.s3Bucket, deps.s3Endpoint)
	if err != nil {
		panic(err)
	}

	applicationRepo := adapters.NewApplicationRepositoryPostgre(deps.pg)

	svc := service.NewService(spreadsheetSvc, applicationRepo, spreadsheetStorage)
	return svc
}

type Option func(*dependencies)

type dependencies struct {
	credentials []byte
	pg          *pgxpool.Pool

	s3AccessKey string
	s3SecretKey string
	s3Endpoint  string
	s3Bucket    string
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

func WithStorageS3(s3AccessKey, s3SecretKey, s3Endpoint, s3Bucket string) Option {
	return func(d *dependencies) {
		d.s3AccessKey = s3AccessKey
		d.s3SecretKey = s3SecretKey
		d.s3Endpoint = s3Endpoint
		d.s3Bucket = s3Bucket
	}
}
