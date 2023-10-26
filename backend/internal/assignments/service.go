package assignments

import (
	"context"

	managerAdapters "github.com/doodocs/qaztrade/backend/internal/manager/adapters"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/doodocs/qaztrade/backend/internal/assignments/adapters"
	"github.com/doodocs/qaztrade/backend/internal/assignments/service"
	"github.com/doodocs/qaztrade/backend/pkg/doodocs"
	"github.com/doodocs/qaztrade/backend/pkg/emailer"
	"github.com/doodocs/qaztrade/backend/pkg/publisher"
	"github.com/doodocs/qaztrade/backend/pkg/spreadsheets"
	"github.com/doodocs/qaztrade/backend/pkg/storage"
	"github.com/jackc/pgx/v4/pgxpool"
)

func MakeService(ctx context.Context, opts ...Option) service.Service {
	deps := &dependencies{}
	deps.setDefaults()
	for _, opt := range opts {
		opt(deps)
	}

	storage, err := storage.NewStorageS3(ctx, deps.s3AccessKey, deps.s3SecretKey, deps.s3Bucket, deps.s3Endpoint)
	if err != nil {
		panic(err)
	}

	spreadsheetsRepo, err := spreadsheets.NewSpreadsheetClient(ctx, deps.credentialsSA)
	if err != nil {
		panic(err)
	}

	var (
		assignmentsRepo = adapters.NewAssignmentsRepositoryPostgres(deps.pg)
		publisher       = publisher.NewPublisherClient(deps.publisher, deps.topicCheckAssignment)
		emailSvc        = emailer.NewEmailerClient(deps.mailEmail, deps.mailPassword)
		msgRepo         = adapters.NewMessagesRepositoryPostgres(deps.pg)
		applicationRepo = managerAdapters.NewApplicationRepositoryPostgres(deps.pg)
		doodocs         = doodocs.NewDoodocsClient(deps.signUrlBase, deps.signLogin, deps.signPassword)
	)

	svc := service.NewService(
		assignmentsRepo,
		storage,
		spreadsheetsRepo,
		publisher,
		emailSvc,
		msgRepo,
		applicationRepo,
		doodocs,
	)
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

	mailEmail, mailPassword string

	signUrlBase  string
	signLogin    string
	signPassword string
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

func WithMail(mailEmail, mailPassword string) Option {
	return func(d *dependencies) {
		d.mailEmail = mailEmail
		d.mailPassword = mailPassword
	}
}

func WithSignCredentials(signUrlBase, signLogin, signPassword string) Option {
	return func(d *dependencies) {
		d.signUrlBase = signUrlBase
		d.signLogin = signLogin
		d.signPassword = signPassword
	}
}
