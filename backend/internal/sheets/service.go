package sheets

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sheets/adapters"
	"github.com/doodocs/qaztrade/backend/internal/sheets/service"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
)

func MakeService(ctx context.Context, opts ...Option) service.Service {
	deps := &dependencies{}
	deps.setDefaults()
	for _, opt := range opts {
		opt(deps)
	}

	sheetsRepo, err := adapters.NewSpreadsheetClient(ctx, deps.credentials)
	if err != nil {
		panic(err)
	}

	storage, err := adapters.NewStorageS3(ctx, deps.s3AccessKey, deps.s3SecretKey, deps.s3Bucket, deps.s3Endpoint)
	if err != nil {
		panic(err)
	}

	var (
		applicationRepo           = adapters.NewApplicationRepositoryPostgre(deps.pg)
		spreadsheetDevMetadataSvc = sheetsRepo.NewSpreadsheetServiceMetadata()
	)
	svc := service.NewService(sheetsRepo, storage, applicationRepo, *spreadsheetDevMetadataSvc)
	return svc
}

type Option func(*dependencies)

type dependencies struct {
	credentials         []byte
	s3AccessKey         string
	s3SecretKey         string
	s3Endpoint          string
	s3Bucket            string
	originSpreadsheetID string
	pg                  *pgxpool.Pool
}

func (d *dependencies) setDefaults() {
	// pass
}

func WithSheetsCredentials(credentials []byte) Option {
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

func WithPostgre(pg *pgxpool.Pool) Option {
	return func(d *dependencies) {
		d.pg = pg
	}
}

func WithOAuthCredentials(clientSecretBytes []byte) Option {
	return func(d *dependencies) {
		d.clientSecretBytes = clientSecretBytes
	}
}

func WithServiceAccount(svcAccount string) Option {
	return func(d *dependencies) {
		d.svcAccount = svcAccount
	}
}

func WithReviewer(reviewerAccount string) Option {
	return func(d *dependencies) {
		d.reviewerAccount = reviewerAccount
	}
}

func WithOriginSpreadsheetID(originSpreadsheetID string) Option {
	return func(d *dependencies) {
		d.originSpreadsheetID = originSpreadsheetID
	}
}

func WithTemplateSpreadsheetID(templateSpreadsheetID string) Option {
	return func(d *dependencies) {
		d.templateSpreadsheetID = templateSpreadsheetID
	}
}

func WithDestinationFolderID(destinationFolderID string) Option {
	return func(d *dependencies) {
		d.destinationFolderID = destinationFolderID
	}
}

func WithJWT(jwtcli *jwt.Client) Option {
	return func(d *dependencies) {
		d.jwtcli = jwtcli
	}
}
