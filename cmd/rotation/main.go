package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/iKOPKACtraxa/project_rotation/internal/app"
	"github.com/iKOPKACtraxa/project_rotation/internal/logger"
	grpcserver "github.com/iKOPKACtraxa/project_rotation/internal/server/grpc"
	sqlstorage "github.com/iKOPKACtraxa/project_rotation/internal/storage/sql"
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "config", "../../configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig()

	logg := logger.New(config.Logger.File, config.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGHUP)
	defer cancel()

	storage, err := sqlstorage.New(ctx, config.Storage.ConnStr, logg)
	if err != nil {
		logg.Error("at sql storage creating has got an error:", err)
	}

	rotation := app.New(logg, storage)

	err = grpcserver.Serve(ctx, rotation, config.GRPCServer.HostPort)
	if err != nil {
		logg.Error("at grpcserver starting has got an error:", err)
	}
}
