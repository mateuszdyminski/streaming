package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
	"github.com/mateuszdyminski/streaming/cfg"
	"github.com/mateuszdyminski/streaming/handlers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Variables injected by -X flag
var appVersion = "unknown"
var lastCommitTime = "unknown"
var lastCommitHash = "unknown"
var lastCommitUser = "unknown"
var buildTime = "unknown"

// Application is the application object that runs HTTP server.
type Application struct {
}

func (app *Application) middlewares(cfg *cfg.Config) (*interpose.Middleware, error) {
	middle := interpose.New()
	router := app.mux(cfg)

	gzipHandler := gziphandler.GzipHandler(handlers.NewLogginHandler(handlers.NewMetricsHandler(router)))
	middle.UseHandler(gzipHandler)

	return middle, nil
}

func (app *Application) mux(cfg *cfg.Config) *mux.Router {
	router := mux.NewRouter()

	router.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{},
	))

	// initialize handlers per group of endpoints
	handlers.ConfigureStreamingRest(cfg, router)
	handlers.ConfigureHealthRest(cfg, router)

	// Path of static files must be last!
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(cfg.StaticsPath)))

	return router
}

func loadConfig() (*cfg.Config, error) {
	conf, err := cfg.LoadCfg()
	if err != nil {
		return nil, err
	}

	conf.Build = cfg.Build{
		Version:   appVersion,
		BuildTime: buildTime,
		LastCommit: cfg.Commit{
			Author: lastCommitUser,
			Id:     lastCommitHash,
			Time:   lastCommitTime,
		},
	}

	conf.Print()

	return conf, nil
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("[fatal] can't laod configuration! err: %v\n", err)
	}

	app := Application{}
	middle, err := app.middlewares(cfg)
	if err != nil {
		log.Fatalf("[fatal] can't create http middlewares! err: %v\n", err)
	}

	drainInterval, err := time.ParseDuration(cfg.HttpDrainInterval)
	if err != nil {
		log.Fatalf("[fatal] can't parse drain interval! err: %v\n", err)
	}

	// subscribe to SIGINT signals
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)

	// create and start http server in new goroutine
	srv := &http.Server{Addr: cfg.Host, Handler: middle}
	go func() {
		// we can't use log.Fatal here!
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("[warn] http server stoped: %s\n", err)
		}
	}()

	log.Printf("[info] running HTTP server on %s\n", cfg.Host)

	// blocks the execution until os.Interrupt or os.Kill signal appears
	<-quit
	log.Println("[info] shutting down server. waiting to drain the ongoing requests...")

	// shut down gracefully, but wait no longer than the Timeout value.
	ctx, _ := context.WithTimeout(context.Background(), drainInterval)

	// shutdown the http server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[fatal] error while shutdown http server: %v\n", err)
	}

	log.Println("[info] server gracefully stopped")
}
