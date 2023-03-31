package jwt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	var (
		tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzaWQiOiIxSTd0WUFoVWpQSkdhTVU3X1hiaEMwOHJRdzU1SVJjN2JFdGcxbWdtUlBLZyJ9.7yAZuGAm7_WSkGJURMSn5aS8UacVAY-CPx-vOO0rPDE"
		secret      = "qaztradesecret"
		expClaims   = &Claims{
			SpreadsheetID: "1I7tYAhUjPJGaMU7_XbhC08rQw55IRc7bEtg1mgmRPKg",
		}

		jwtcli = NewClient(secret)
	)

	claims, err := jwtcli.Parse(tokenString)
	require.Nil(t, err)
	require.Equal(t, expClaims, claims)
}

func TestNewTokenString(t *testing.T) {
	var (
		token     = "1I7tYAhUjPJGaMU7_XbhC08rQw55IRc7bEtg1mgmRPKg"
		secret    = "qaztradesecret"
		expClaims = &Claims{
			SpreadsheetID: "1I7tYAhUjPJGaMU7_XbhC08rQw55IRc7bEtg1mgmRPKg",
		}

		jwtcli = NewClient(secret)
	)

	tokenString, err := jwtcli.NewTokenString(token)
	require.Nil(t, err)
	require.NotEmpty(t, tokenString)

	claims, err := jwtcli.Parse(tokenString)
	require.Nil(t, err)
	require.Equal(t, expClaims, claims)
}
