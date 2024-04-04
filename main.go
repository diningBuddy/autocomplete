package main

import (
	"context"
	"github.com/skku/autocomplete/pkg/config"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/skku/autocomplete/app"
	"github.com/skku/autocomplete/env"

	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

var properties *config.Properties

func init() {
	configPath := env.GetString("CONFIG_PATH", "./example-config.yaml")
	if c, err := config.LoadConfig(configPath); err != nil {
		log.Fatalf("server booting failed: %s", err)
	} else {
		properties = c
	}
	log.SetFormatter(
		&easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%] %time% - %msg%\n",
		},
	)
}

func main() {
	app := &app.App{}
	app.Initialize(properties)
	go app.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server...")
	timeout := 5 * time.Second
	log.Infof("shutting down server after 5 sec...")
	time.Sleep(timeout)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	tracer.Stop()

	if err := app.GracefulShutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown:", err)
	}
}
