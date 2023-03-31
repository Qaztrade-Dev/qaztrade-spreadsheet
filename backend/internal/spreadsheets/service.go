package spreadsheets

import (
	"context"

	"github.com/doodocs/qaztrade/backend/internal/sheets/service"
)

func MakeService(ctx context.Context, opts ...Option) service.Service {
	deps := &dependencies{}
	deps.setDefaults()
	for _, opt := range opts {
		opt(deps)
	}

	// sheetsRepo, err := adapters.NewSpreadsheetClient(ctx, deps.credentials)
	// if err != nil {
	// 	panic(err)
	// }

	// storage, err := adapters.NewStorageS3(ctx, deps.s3AccessKey, deps.s3SecretKey, deps.s3Bucket, deps.s3Endpoint)
	// if err != nil {
	// 	panic(err)
	// }

	// svc := service.NewService(sheetsRepo, storage)
	// return svc
	return nil
}

type Option func(*dependencies)

type dependencies struct {
}

func (d *dependencies) setDefaults() {
	// pass
}

// func WithSheetsCredentials(credentials []byte) Option {
// 	return func(d *dependencies) {
// 		d.credentials = credentials
// 	}
// }