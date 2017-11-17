package main

import (
	"log"
	"net/http"
	"time"

	"github.com/carbocation/interpose"
	"github.com/gorilla/mux"
	"github.com/mateuszdyminski/streaming/cfg"
	"github.com/mateuszdyminski/streaming/handlers"
	"github.com/tylerb/graceful"
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
	middle.Use(handlers.LogRequest())
	middle.UseHandler(app.mux(cfg))

	return middle, nil
}

func (app *Application) mux(cfg *cfg.Config) *mux.Router {
	router := mux.NewRouter()

	// initialize handlers per group of endpoints
	handlers.ConfigureStreamingRest(cfg, router)
	handlers.ConfigureHealthRest(cfg, router)

	// Path of static files must be last!
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(cfg.StaticsPath)))

	return router
}

func NewApp(cfg *cfg.Config) (*Application, error) {
	return &Application{}, nil
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
	if err != nil {
		log.Fatalf("[fatal] can't create application! err: %v\n", err)
	}

	middle, err := app.middlewares(cfg)
	if err != nil {
		log.Fatalf("[fatal] can't create http middlewares! err: %v\n", err)
	}

	drainInterval, err := time.ParseDuration(cfg.HttpDrainInterval)
	if err != nil {
		log.Fatalf("[fatal] can't parse drain interval! err: %v\n", err)
	}

	srv := &graceful.Server{
		Timeout: drainInterval,
		Server:  &http.Server{Addr: cfg.Host, Handler: middle},
	}

	log.Printf("[info] running HTTP server on %s\n", cfg.Host)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
