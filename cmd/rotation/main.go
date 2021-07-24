package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/iKOPKACtraxa/otus-hw/project_rotation/internal/app"
	"github.com/iKOPKACtraxa/otus-hw/project_rotation/internal/logger"
	sqlstorage "github.com/iKOPKACtraxa/otus-hw/project_rotation/internal/storage/sql"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	storage, err := sqlstorage.New(ctx, config.Storage.ConnStr, logg)
	if err != nil {
		logg.Error("at sql storage creating has got an error: " + err.Error())
	}
	rotation := app.New(logg, storage) // todo1
	// server := internalHTTP.NewServer(rotation, config.HTTPServer.HostPort) // тут будет GRPC сервер см. main из elections и еще конфиг нужно подшаманить
	_ = rotation // todo1
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP)

		select {
		case <-ctx.Done():
			return
		case <-signals:
		}

		signal.Stop(signals)
		cancel()
		/*todo1
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {//todo1 тут какуюто остановку grpc сервера если это нунжно (в elections вроде не было)
			logg.Error("failed to stop http server: " + err.Error())
		}
		*/
	}()

	logg.Info("rotation is running...")

	// для теста
	// err = rotation.DeleteBanner(ctx, 666, 1)
	// logg.Error("" + err.Error())

	/*todo1
	if err := server.Start(ctx); err != nil { // это вроде старт самого сервера (опять же см elections)
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
	*/
}
