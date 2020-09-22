package main

import (
	"os"
	"os/signal"
	"syscall"

	"projecto/app"
	"projecto/config"
	"projecto/web"

	log "github.com/sirupsen/logrus"
)

func main() {
	logger := log.Fields{"app": true}
	app := app.App{Logger: log.WithFields(logger)}
	initialize(&app)
	if err := app.Start(); err != nil {
		log.WithFields(logger).Fatalf("app %v %v cant start error: %v", app.Name(), app.Version(), err)
	}
	shutdownsignals := make(chan os.Signal, 1)
	signal.Notify(shutdownsignals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-shutdownsignals
	log.WithFields(logger).Errorf("app %v %v received signal: %v", app.Name(), app.Version(), sig)
	if err := app.Close(); err != nil {
		log.WithFields(logger).Errorf("cant stop app: %v", err)
	}
}

func initialize(a *app.App) {
	a.Register(config.New()).
		Register(web.New())
}
