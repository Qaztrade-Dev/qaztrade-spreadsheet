package jwt

import (
	"os"
	"testing"
	"time"

	stdjwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

type testClaims struct {
	SpreadsheetID string `json:"sid"`
}

func TestParse(t *testing.T) {
	var (
		tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwIjp7InNpZCI6IjFJN3RZQWhValBKR2FNVTdfWGJoQzA4clF3NTVJUmM3YkV0ZzFtZ21SUEtnIn19.QfFWe6asN8YVA-hElikkOAq1jlueW3_e-9oohf_xf0k"
		secret      = "qaztradesecret"
		expClaims   = &testClaims{
			SpreadsheetID: os.Getenv("TEMPLATE_SPREADSHEET_ID"),
		}

		jwtcli = NewClient(secret)
	)

	claims, err := Parse[testClaims](jwtcli, tokenString)
	require.Nil(t, err)
	require.Equal(t, expClaims, claims)
}

func TestNewTokenString(t *testing.T) {
	var (
		spreadhsheetID = os.Getenv("TEMPLATE_SPREADSHEET_ID")
		secret         = "qaztradesecret"
		expClaims      = &testClaims{
			SpreadsheetID: os.Getenv("TEMPLATE_SPREADSHEET_ID"),
		}

		jwtcli = NewClient(secret)
	)

	tokenString, err := NewTokenString(jwtcli, &testClaims{
		SpreadsheetID: spreadhsheetID,
	})
	require.Nil(t, err)
	require.NotEmpty(t, tokenString)

	claims, err := Parse[testClaims](jwtcli, tokenString)
	require.Nil(t, err)
	require.Equal(t, expClaims, claims)
}

func TestWithExpire(t *testing.T) {
	var (
		spreadhsheetID = os.Getenv("TEMPLATE_SPREADSHEET_ID")
		secret         = "qaztradesecret"

		jwtcli = NewClient(secret)
	)

	tokenString, err := NewTokenString(
		jwtcli,
		&testClaims{SpreadsheetID: spreadhsheetID},
		WithExpire(time.Date(2019, 3, 30, 1, 2, 3, 4, time.UTC)),
	)
	require.Nil(t, err)
	require.NotEmpty(t, tokenString)

	claims, err := Parse[testClaims](jwtcli, tokenString)
	require.ErrorIs(t, err, stdjwt.ErrTokenExpired)
	require.Equal(t, (*testClaims)(nil), claims)
}
