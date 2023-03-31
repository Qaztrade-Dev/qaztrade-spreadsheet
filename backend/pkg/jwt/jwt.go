package jwt

import (
	"time"

	stdjwt "github.com/golang-jwt/jwt/v5"
)

type Client struct {
	secret []byte
}

func NewClient(secret string) *Client {
	return &Client{
		secret: []byte(secret),
	}
}

func Parse[T any](c *Client, tokenString string) (*T, error) {
	type tmp struct {
		Payload *T `json:"p"`
		stdjwt.RegisteredClaims
	}

	token, err := stdjwt.ParseWithClaims(tokenString, &tmp{}, func(token *stdjwt.Token) (interface{}, error) {
		return c.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*tmp); ok && token.Valid {
		return claims.Payload, nil
	}

	return nil, err
}

func NewTokenString[T any](c *Client, t *T, opts ...Option) (string, error) {
	type tmp struct {
		Payload *T `json:"p"`
		stdjwt.RegisteredClaims
	}

	claims := tmp{Payload: t}
	for _, opt := range opts {
		opt(&claims.RegisteredClaims)
	}

	token := stdjwt.NewWithClaims(stdjwt.SigningMethodHS256, claims)

	return token.SignedString(c.secret)
}

type Option func(claims *stdjwt.RegisteredClaims)

func WithExpire(t time.Time) Option {
	return func(claims *stdjwt.RegisteredClaims) {
		claims.ExpiresAt = stdjwt.NewNumericDate(t)
	}
}
