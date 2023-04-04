package spreadsheets

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/adapters"
	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/service"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
)

func MakeService(ctx context.Context, opts ...Option) service.Service {
	deps := &dependencies{}
	deps.setDefaults()
	for _, opt := range opts {
		opt(deps)
	}

	spreadsheetSvc, err := adapters.NewSpreadsheetServiceGoogle(
		deps.clientSecretBytes,
		deps.svcAccount,
		deps.reviewerAccount,
		deps.jwtcli,
		deps.pg,
		deps.templateSpreadsheetID,
		deps.destinationFolderID,
	)
	if err != nil {
		panic(err)
	}

	var (
		applicationRepo = adapters.NewApplicationRepositoryPostgre(deps.pg)
		userRepo        = adapters.NewUserRepositoryPostgre(deps.pg)
	)

	svc := service.NewService(spreadsheetSvc, applicationRepo, userRepo)
	return svc
}

type Option func(*dependencies)

type dependencies struct {
	clientSecretBytes     []byte
	svcAccount            string
	reviewerAccount       string
	jwtcli                *jwt.Client
	pg                    *pgxpool.Pool
	templateSpreadsheetID string
	destinationFolderID   string
}

func (d *dependencies) setDefaults() {
	// pass
}

func WithJWT(jwtcli *jwt.Client) Option {
	return func(d *dependencies) {
		d.jwtcli = jwtcli
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
