package adapters

import (
	"context"
	"testing"

	"github.com/doodocs/qaztrade/backend/internal/auth/domain"
	"github.com/doodocs/qaztrade/backend/internal/auth/service"
	"github.com/stretchr/testify/require"
)

func TestEmailSvcSend(t *testing.T) {
	var (
		ctx      = context.Background()
		email    = "qaztrade.export@gmail.com"
		password = "oqbenkgtnddsaqpl"
		svc      = NewEmailServiceGmail(email, password)
	)

	err := svc.Send(ctx, "ali.tlekbai+test@gmail.com", "", &service.MailPayload{
		Credentials: &domain.Credentials{
			AccessToken: "1",
		},
	})
	require.Nil(t, err)
}
