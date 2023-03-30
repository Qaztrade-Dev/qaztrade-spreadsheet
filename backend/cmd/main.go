package main

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/doodocs/qaztrade/backend/internal/common"
	"github.com/doodocs/qaztrade/backend/internal/sheets"
	"github.com/doodocs/qaztrade/backend/internal/sheets/pkg/jwt"
	"github.com/go-kit/log"
)

//go:embed credentials.json
var credentials []byte

const (
	defaultPort = "8082"
)

func main() {
	var (
		ctx         = context.Background()
		port        = getenv("PORT", defaultPort)
		jwtsecret   = getenv("JWT_SECRET", "qaztradesecret")
		s3AccessKey = getenv("S3_ACCESS_KEY", "")
		s3SecretKey = getenv("S3_SECRET_KEY", "")
		s3Endpoint  = getenv("S3_ENDPOINT", "")
		s3Bucket    = getenv("S3_BUCKET", "")
		addr        = ":" + port
	)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	}

	var (
		sheetsService = sheets.MakeService(
			ctx,
			sheets.WithSheetsCredentials(credentials),
			sheets.WithStorageS3(s3AccessKey, s3SecretKey, s3Endpoint, s3Bucket),
		)
	)

	var (
		jwtcli     = jwt.NewClient(jwtsecret)
		httpLogger = log.With(logger, "component", "http")
		mux        = http.NewServeMux()
	)

	mux.Handle("/sheets/", sheets.MakeHandler(sheetsService, jwtcli, httpLogger))

	http.Handle("/", common.AccessControl(mux))

	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", addr, "msg", "listening")
		errs <- http.ListenAndServe(addr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func getenv(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
