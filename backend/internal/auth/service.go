package auth

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/auth/service"
)

func MakeService(ctx context.Context, opts ...Option) service.Service {
	deps := &dependencies{}
	deps.setDefaults()
	for _, opt := range opts {
		opt(deps)
	}

	// svc := service.NewService(sheetsRepo, storage)
	return nil
}

type Option func(*dependencies)

type dependencies struct {
	// pass
}

func (d *dependencies) setDefaults() {
	// pass
}

func WithSheetsCredentials(credentials []byte) Option {
	return func(d *dependencies) {
		// d.credentials = credentials
	}
}
