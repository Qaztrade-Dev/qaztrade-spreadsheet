package assignments

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
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

	storage, err := adapters.NewStorageS3(ctx, deps.s3AccessKey, deps.s3SecretKey, deps.s3Bucket, deps.s3Endpoint)
	if err != nil {
		panic(err)
	}

	spreadsheetsRepo, err := adapters.NewSpreadsheetClient(ctx, deps.credentialsSA)
	if err != nil {
		panic(err)
	}

	var (
		assignmentsRepo = adapters.NewAssignmentsRepositoryPostgres(deps.pg)
		publisher       = adapters.NewPublisherWatermill(deps.publisher, deps.topicCheckAssignment)
	)

	svc := service.NewService(assignmentsRepo, storage, spreadsheetsRepo, publisher)
	return svc
}

type Option func(*dependencies)

type dependencies struct {
	pg            *pgxpool.Pool
	credentialsSA []byte

	s3AccessKey string
	s3SecretKey string
	s3Endpoint  string
	s3Bucket    string

	publisher            message.Publisher
	topicCheckAssignment string
}

func (d *dependencies) setDefaults() {
	// pass
}

func WithPostgres(pg *pgxpool.Pool) Option {
	return func(d *dependencies) {
		d.pg = pg
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

func WithCredentialsSA(credentialsSA []byte) Option {
	return func(d *dependencies) {
		d.credentialsSA = credentialsSA
	}
}

func WithPublisher(publisher message.Publisher, topicCheckAssignment string) Option {
	return func(d *dependencies) {
		d.publisher = publisher
		d.topicCheckAssignment = topicCheckAssignment
	}
}
