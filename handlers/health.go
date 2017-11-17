package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/mateuszdyminski/streaming/cfg"
)

type HealthApi struct {
	cfg       *cfg.Config
	hostname  string
	startedAt time.Time
}

func ConfigureHealthRest(cfg *cfg.Config, r *mux.Router) {
	host, err := os.Hostname()
	if err != nil {
		host = "localhost"
	}

	rest := &HealthApi{
		cfg:       cfg,
		hostname:  host,
		startedAt: time.Now().UTC(),
	}

	// health API
	r.HandleFunc("/api/health", rest.health).Methods("GET")

}

type Response struct {
	Build     cfg.Build `json:"build"`
	Hostname  string    `json:"hostname"`
	Uptime    string    `json:"uptime"`
	StartedAt string    `json:"startedAt"`
}

func (a *HealthApi) health(w http.ResponseWriter, req *http.Request) {
	resp := Response{
		Hostname:  a.hostname,
		StartedAt: a.startedAt.Format("2006-01-02_15:04:05"),
		Uptime:    time.Now().UTC().Sub(a.startedAt).String(),
		Build:     a.cfg.Build,
	}

	if err := WriteJSON(w, resp); err != nil {
		WriteErr(w, err, http.StatusInternalServerError)
		return
	}
}
