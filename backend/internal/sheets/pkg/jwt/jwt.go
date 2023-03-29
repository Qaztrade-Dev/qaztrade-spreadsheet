package jwt

import (
	stdjwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	SpreadsheetID string `json:"sid"`
	stdjwt.RegisteredClaims
}

type Client struct {
	secret []byte
}

func NewClient(secret string) *Client {
	return &Client{
		secret: []byte(secret),
	}
}

func (c *Client) Parse(tokenString string) (*Claims, error) {
	token, err := stdjwt.ParseWithClaims(tokenString, &Claims{}, func(token *stdjwt.Token) (interface{}, error) {
		return c.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func (c *Client) NewTokenString(spreadsheetID string) (string, error) {
	var (
		claims = Claims{SpreadsheetID: spreadsheetID}
		token  = stdjwt.NewWithClaims(stdjwt.SigningMethodHS256, claims)
	)

	return token.SignedString(c.secret)
}
