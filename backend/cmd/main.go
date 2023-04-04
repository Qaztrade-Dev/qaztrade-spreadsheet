package main

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/doodocs/qaztrade/backend/internal/auth"
	"github.com/doodocs/qaztrade/backend/internal/common"
	"github.com/doodocs/qaztrade/backend/internal/google"
	"github.com/doodocs/qaztrade/backend/internal/manager"
	"github.com/doodocs/qaztrade/backend/internal/sheets"
	"github.com/doodocs/qaztrade/backend/internal/spreadsheets"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/go-kit/log"
	"github.com/jackc/pgx/v4/pgxpool"
	googleOAuth "golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	googleSheets "google.golang.org/api/sheets/v4"
)

//go:embed credentials.json
var credentials []byte

//go:embed oauth_secret.json
var oauthSecret []byte

const (
	defaultPort = "8082"
)

func main() {
	var (
		ctx                   = context.Background()
		port                  = getenv("PORT", defaultPort)
		jwtsecret             = getenv("JWT_SECRET", "qaztradesecret")
		s3AccessKey           = getenv("S3_ACCESS_KEY")
		s3SecretKey           = getenv("S3_SECRET_KEY")
		s3Endpoint            = getenv("S3_ENDPOINT")
		s3Bucket              = getenv("S3_BUCKET")
		postgresLogin         = getenv("POSTGRES_LOGIN", "postgres")
		postgresPassword      = getenv("POSTGRES_PASSWORD", "postgres")
		postgresHost          = getenv("POSTGRES_HOST", "localhost")
		postgresDatabase      = getenv("POSTGRES_DATABASE", "qaztrade")
		mailLogin             = getenv("MAIL_LOGIN")
		mailPassword          = getenv("MAIL_PASSWORD")
		svcAccount            = getenv("SERVICE_ACCOUNT")
		templateSpreadsheetID = getenv("TEMPLATE_SPREADSHEET_ID")
		destinationFolderID   = getenv("DESTINATION_FOLDER_ID")
		originSpreadsheetID   = getenv("ORIGIN_SPREADSHEET_ID")

		addr        = ":" + port
		postgresURL = fmt.Sprintf("postgresql://%s:%s@%s:5432/%s", postgresLogin, postgresPassword, postgresHost, postgresDatabase)
		jwtcli      = jwt.NewClient(jwtsecret)
	)

	pg, err := pgxpool.Connect(ctx, postgresURL)
	if err != nil {
		panic(err)
	}

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	}

	oauthConfig, err := googleOAuth.ConfigFromJSON(oauthSecret, drive.DriveScope, googleSheets.SpreadsheetsScope)
	if err != nil {
		panic(err)
	}

	var (
		sheetsService = sheets.MakeService(
			ctx,
			sheets.WithSheetsCredentials(credentials),
			sheets.WithStorageS3(s3AccessKey, s3SecretKey, s3Endpoint, s3Bucket),
			sheets.WithOriginSpreadsheetID(originSpreadsheetID),
		)

		authService = auth.MakeService(
			ctx,
			auth.WithPostgre(pg),
			auth.WithJWT(jwtcli),
			auth.WithMail(mailLogin, mailPassword),
		)

		googleService = google.MakeService(
			ctx,
			google.WithPostgre(pg),
			google.WithOAuthConfig(oauthConfig),
		)

		spreadsheetsService = spreadsheets.MakeService(
			ctx,
			spreadsheets.WithPostgre(pg),
			spreadsheets.WithJWT(jwtcli),
			spreadsheets.WithOAuthCredentials(oauthSecret),
			spreadsheets.WithServiceAccount(svcAccount),
			spreadsheets.WithTemplateSpreadsheetID(templateSpreadsheetID),
			spreadsheets.WithDestinationFolderID(destinationFolderID),
		)

		managerService = manager.MakeService(
			ctx,
			manager.WithPostgre(pg),
			manager.WithCredentials(credentials),
		)
	)

	var (
		httpLogger = log.With(logger, "component", "http")
		mux        = http.NewServeMux()
	)

	mux.Handle("/sheets/", sheets.MakeHandler(sheetsService, jwtcli, httpLogger))
	mux.Handle("/auth/", auth.MakeHandler(authService, jwtcli, httpLogger))
	mux.Handle("/google/", google.MakeHandler(googleService, httpLogger))
	mux.Handle("/spreadsheets/", spreadsheets.MakeHandler(spreadsheetsService, jwtcli, httpLogger))
	mux.Handle("/manager/", manager.MakeHandler(managerService, jwtcli, httpLogger))

	http.Handle("/", common.AccessControl(mux))

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", addr, "msg", "listening")
		errs <- http.ListenAndServe(addr, nil)
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func getenv(env string, fallback ...string) string {
	e := os.Getenv(env)
	if e == "" {
		value := ""
		if len(fallback) > 0 {
			value = fallback[0]
		}
		return value
	}
	return e
}
