package qaztradeoauth2

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/sheets/v4"
)

type Client struct {
	pg     *pgxpool.Pool
	config *oauth2.Config
	mu     *sync.Mutex
}

func NewClient(credentialsOAuth []byte, pg *pgxpool.Pool) (*Client, error) {
	config, err := google.ConfigFromJSON(credentialsOAuth, drive.DriveScope, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}

	result := &Client{
		config: config,
		pg:     pg,
		mu:     &sync.Mutex{},
	}

	return result, nil
}

type customTokenSource struct {
	oauth2.TokenSource
	mu             *sync.Mutex
	InitialToken   *oauth2.Token
	qaztradeOAuth2 *Client
	ctx            context.Context
}

func (cts *customTokenSource) Token() (*oauth2.Token, error) {
	cts.mu.Lock()
	defer cts.mu.Unlock()

	token, err := cts.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	if cts.InitialToken.AccessToken != token.AccessToken || cts.InitialToken.RefreshToken != token.RefreshToken {
		if err := cts.qaztradeOAuth2.updateToken(cts.ctx, token); err != nil {
			return nil, err
		}
	}

	return token, nil
}

func (o *Client) GetClient(ctx context.Context) (*http.Client, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	tokenStr, err := o.getOauthToken(ctx)
	if err != nil {
		return nil, err
	}

	tokenOauth2, err := tokenFromStr(tokenStr)
	if err != nil {
		return nil, err
	}

	tokenSource := o.config.TokenSource(ctx, tokenOauth2)
	customTS := &customTokenSource{
		InitialToken:   tokenOauth2,
		TokenSource:    oauth2.ReuseTokenSource(tokenOauth2, tokenSource),
		qaztradeOAuth2: o,
		mu:             o.mu,
		ctx:            ctx,
	}

	// httpCli := clientWithToken(s.config, tokenOauth2)
	client := oauth2.NewClient(ctx, customTS)

	return client, nil
}

func (o *Client) getOauthToken(ctx context.Context) (string, error) {
	const sql = `
		select 
			token
		from "oauth2_tokens"
		where id = 1
	`

	var (
		token string
	)

	err := o.pg.QueryRow(ctx, sql).Scan(&token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func tokenFromStr(tokenStr string) (*oauth2.Token, error) {
	tok := &oauth2.Token{}
	err := json.Unmarshal([]byte(tokenStr), &tok)
	return tok, err
}

// func clientWithToken(config *oauth2.Config, token *oauth2.Token) *http.Client {
// 	return config.Client(context.Background(), token)
// }

func (o *Client) updateToken(ctx context.Context, token *oauth2.Token) error {
	tokenStr, err := encodeToken(token)
	if err != nil {
		return err
	}

	const sql = `
		insert into "oauth2_tokens" ("id", "token") values 
			(1, $1)
		on conflict ("id")
		do update set "token" = excluded."token"
	`

	if _, err := o.pg.Exec(ctx, sql, tokenStr); err != nil {
		return err
	}

	return nil
}

func encodeToken(token *oauth2.Token) (string, error) {
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return "", err
	}

	return string(tokenBytes), nil
}
