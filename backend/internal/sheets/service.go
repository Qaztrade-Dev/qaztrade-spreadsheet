package sheets

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sheets/adapters"
	"github.com/doodocs/qaztrade/backend/internal/sheets/service"
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

	svc := service.NewService(sheetsRepo)
	return svc
}

type Option func(*dependencies)

type dependencies struct {
	credentials []byte
}

func (d *dependencies) setDefaults() {
	// pass
}

func WithSheetsCredentials(credentials []byte) Option {
	return func(d *dependencies) {
		d.credentials = credentials
	}
}
