package spreadsheets

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/adapters"
	"github.com/doodocs/qaztrade/backend/internal/spreadsheets/service"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/doodocs/qaztrade/backend/pkg/qaztradeoauth2"
	"github.com/jackc/pgx/v4/pgxpool"
)

func MakeService(ctx context.Context, opts ...Option) service.Service {
	deps := &dependencies{}
	deps.setDefaults()
	for _, opt := range opts {
		opt(deps)
	}

	oauth2, err := qaztradeoauth2.NewClient(deps.clientSecretBytes, deps.pg)
	if err != nil {
		panic(err)
	}

	var (
		spreadsheetSvc = adapters.NewSpreadsheetServiceGoogle(
			oauth2,
			deps.svcAccount,
			deps.reviewerAccount,
			deps.jwtcli,
			deps.originSpreadsheetID,
			deps.templateSpreadsheetID,
			deps.destinationFolderID,
		)
		applicationRepo = adapters.NewApplicationRepositoryPostgre(deps.pg)
		userRepo        = adapters.NewUserRepositoryPostgre(deps.pg)
	)

	svc := service.NewService(spreadsheetSvc, applicationRepo, userRepo, deps.jwtcli)
	return svc
}

type Option func(*dependencies)

type dependencies struct {
	clientSecretBytes     []byte
	svcAccount            string
	reviewerAccount       string
	jwtcli                *jwt.Client
	pg                    *pgxpool.Pool
	originSpreadsheetID   string
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
