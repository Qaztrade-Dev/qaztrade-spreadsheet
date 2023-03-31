package jwt

import (
	"testing"

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
			SpreadsheetID: "1I7tYAhUjPJGaMU7_XbhC08rQw55IRc7bEtg1mgmRPKg",
		}

		jwtcli = NewClient(secret)
	)

	claims, err := Parse[testClaims](jwtcli, tokenString)
	require.Nil(t, err)
	require.Equal(t, expClaims, claims)
}

func TestNewTokenString(t *testing.T) {
	var (
		spreadhsheetID = "1I7tYAhUjPJGaMU7_XbhC08rQw55IRc7bEtg1mgmRPKg"
		secret         = "qaztradesecret"
		expClaims      = &testClaims{
			SpreadsheetID: "1I7tYAhUjPJGaMU7_XbhC08rQw55IRc7bEtg1mgmRPKg",
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
