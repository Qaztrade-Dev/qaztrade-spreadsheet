package sign

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/doodocs/qaztrade/backend/internal/sign/adapters"
	"github.com/doodocs/qaztrade/backend/internal/sign/adapters/pdfservice"
	"github.com/doodocs/qaztrade/backend/internal/sign/service"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/doodocs/qaztrade/backend/pkg/publisher"
	"github.com/jackc/pgx/v4/pgxpool"
)

func MakeService(ctx context.Context, opts ...Option) service.Service {
	deps := &dependencies{}
	deps.setDefaults()
	for _, opt := range opts {
		opt(deps)
	}

	pdfSvc, err := pdfservice.NewPDFService()
	if err != nil {
		panic(err)
	}

	spreadsheetSvc, err := adapters.NewSpreadsheetClient(ctx, deps.credentialsSA, deps.adminAccount, deps.svcAccount)
	if err != nil {
		panic(err)
	}

	var (
		signSvc         = adapters.NewSigningServiceDoodocs(deps.signUrlBase, deps.signLogin, deps.signPassword)
		applicationRepo = adapters.NewApplicationRepositoryPostgre(deps.pg)
		publisher       = publisher.NewPublisherClient(deps.publisher, deps.topicDoodocsDocumentSigned)
	)

	svc := service.NewService(pdfSvc, signSvc, spreadsheetSvc, applicationRepo, publisher)
	return svc
}

type Option func(*dependencies)

type dependencies struct {
	credentialsSA []byte
	pg            *pgxpool.Pool
	signUrlBase   string
	signLogin     string
	signPassword  string
	jwtcli        *jwt.Client

	adminAccount string
	svcAccount   string

	publisher                  message.Publisher
	topicDoodocsDocumentSigned string
}

func (d *dependencies) setDefaults() {
	// pass
}

func WithCredentialsSA(credentialsSA []byte) Option {
	return func(d *dependencies) {
		d.credentialsSA = credentialsSA
	}
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

func WithServiceAccount(input string) Option {
	return func(d *dependencies) {
		d.svcAccount = input
	}
}

func WithPublisher(publisher message.Publisher) Option {
	return func(d *dependencies) {
		d.publisher = publisher
	}
}

func WithTopicDoodocsDocumentSigned(input string) Option {
	return func(d *dependencies) {
		d.topicDoodocsDocumentSigned = input
	}
}
