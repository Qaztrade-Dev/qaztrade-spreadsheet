package emailservice

import (
	"bufio"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendNotice(t *testing.T) {
	var (
		ctx      = context.Background()
		email    = "export@p-s.kz"
		password = ""
		svc      = NewEmailServiceGmail(email, password)
	)
	file, err := os.Open("./text_mailname.txt")
	require.Nil(t, err)
	defer file.Close()

	err = svc.SendNotice(ctx, "daniyar.kuttymbek@gmail.com", "Тестирования уведомления", "Yupio.txt", bufio.NewReader(file))
	require.Nil(t, err)

}
