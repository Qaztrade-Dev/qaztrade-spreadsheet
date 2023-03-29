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
		ctx       = context.Background()
		port      = getenv("PORT", defaultPort)
		jwtsecret = getenv("JWT_SECRET", "qaztradesecret")
		addr      = ":" + port
	)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	}

	var (
		sheetsService = sheets.MakeService(ctx, sheets.WithSheetsCredentials(credentials))
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
