package main

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/doodocs/qaztrade/backend/internal/assignments"
	"github.com/doodocs/qaztrade/backend/internal/auth"
	"github.com/doodocs/qaztrade/backend/internal/common"
	"github.com/doodocs/qaztrade/backend/internal/google"
	"github.com/doodocs/qaztrade/backend/internal/manager"
	"github.com/doodocs/qaztrade/backend/internal/sheets"
	"github.com/doodocs/qaztrade/backend/internal/sign"
	"github.com/doodocs/qaztrade/backend/internal/spreadsheets"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/go-kit/log"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/redis/go-redis/v9"
	googleOAuth "golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	googleSheets "google.golang.org/api/sheets/v4"
)

//go:embed credentials_sa.json
var credentialsSA []byte

//go:embed credentials_oauth.json
var credentialsOAuth []byte

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
		reviewerAccount       = getenv("REVIEWER_ACCOUNT")
		adminAccount          = getenv("ADMIN_ACCOUNT")
		templateSpreadsheetID = getenv("TEMPLATE_SPREADSHEET_ID")
		destinationFolderID   = getenv("DESTINATION_FOLDER_ID")
		originSpreadsheetID   = getenv("ORIGIN_SPREADSHEET_ID")
		signUrlBase           = getenv("SIGN_URL_BASE")
		signLogin             = getenv("SIGN_LOGIN")
		signPassword          = getenv("SIGN_PASSWORD")
		redisAddr             = getenv("REDIS_ADDR")
		topicCheckAssignments = getenv("TOPIC_CHECK_ASSIGNMENTS", "check-assignments")

		addr        = ":" + port
		postgresURL = fmt.Sprintf("postgresql://%s:%s@%s:5432/%s", postgresLogin, postgresPassword, postgresHost, postgresDatabase)
		jwtcli      = jwt.NewClient(jwtsecret)
	)

	pg, err := pgxpool.Connect(ctx, postgresURL)
	if err != nil {
		panic(err)
	}

	var (
		logger          log.Logger
		watermillLogger = watermill.NewStdLogger(false, false)
	)

	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	}

	oauthConfig, err := googleOAuth.ConfigFromJSON(credentialsOAuth, drive.DriveScope, googleSheets.SpreadsheetsScope)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	sub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, watermillLogger)
	if err != nil {
		panic(err)
	}

	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, watermillLogger)
	if err != nil {
		panic(err)
	}

	router, err := message.NewRouter(message.RouterConfig{}, watermillLogger)
	if err != nil {
		panic(err)
	}

	var (
		sheetsService = sheets.MakeService(
			ctx,
			sheets.WithSheetsCredentials(credentialsSA),
			sheets.WithPostgre(pg),
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
			spreadsheets.WithOAuthCredentials(credentialsOAuth),
			spreadsheets.WithServiceAccount(svcAccount),
			spreadsheets.WithReviewer(reviewerAccount),
			spreadsheets.WithTemplateSpreadsheetID(templateSpreadsheetID),
			spreadsheets.WithDestinationFolderID(destinationFolderID),
			spreadsheets.WithOriginSpreadsheetID(originSpreadsheetID),
		)

		managerService = manager.MakeService(
			ctx,
			manager.WithPostgre(pg),
			manager.WithCredentials(credentialsSA),
			manager.WithSignCredentials(signUrlBase, signLogin, signPassword),
			manager.WithAdmin(adminAccount),
			manager.WithServiceAccount(svcAccount),
			manager.WithStorageS3(s3AccessKey, s3SecretKey, s3Endpoint, s3Bucket),
			manager.WithMail(mailLogin, mailPassword),
		)

		signService = sign.MakeService(
			ctx,
			sign.WithJWT(jwtcli),
			sign.WithPostgre(pg),
			sign.WithSignCredentials(signUrlBase, signLogin, signPassword),
			sign.WithCredentialsSA(credentialsSA),
			sign.WithAdmin(adminAccount),
			sign.WithServiceAccount(svcAccount),
		)

		assignmentsService = assignments.MakeService(
			ctx,
			assignments.WithPostgres(pg),
			assignments.WithStorageS3(s3AccessKey, s3SecretKey, s3Endpoint, s3Bucket),
			assignments.WithCredentialsSA(credentialsSA),
			assignments.WithPublisher(publisher, topicCheckAssignments),
		)
	)

	router.AddNoPublisherHandler(
		topicCheckAssignments,
		topicCheckAssignments,
		sub,
		func(msg *message.Message) error {
			var (
				assignmentIDStr = string(msg.Payload)
				assignmentID, _ = strconv.ParseUint(assignmentIDStr, 10, 0)
			)

			fmt.Printf("processing - %d\n", assignmentID)
			if err := assignmentsService.CheckAssignment(msg.Context(), assignmentID); err != nil {
				fmt.Println("error happened", err)
				return err
			}

			return nil
		},
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
	mux.Handle("/sign/", sign.MakeHandler(signService, jwtcli, httpLogger))
	mux.Handle("/assignments/", assignments.MakeHandler(assignmentsService, jwtcli, httpLogger))

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

	go func() {
		errs <- router.Run(context.Background())
	}()

	logger.Log("terminated", <-errs)
}

func getenv(env string, fallback ...string) string {
	value := os.Getenv(env)
	if value != "" {
		return value
	}

	if len(fallback) > 0 {
		value = fallback[0]
	}
	return value
}
