package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mateuszdyminski/streaming/cfg"
	"github.com/sirupsen/logrus"
)

type StreamingApi struct {
	cfg *cfg.Config
}

func ConfigureStreamingRest(cfg *cfg.Config, r *mux.Router) {
	rest := &StreamingApi{cfg: cfg}

	// videos API
	r.HandleFunc("/api/videos/{id}", rest.video).Methods("GET")
	r.HandleFunc("/api/posters/{id}", rest.poster).Methods("GET")

}

func (l *StreamingApi) video(w http.ResponseWriter, req *http.Request) {
	videoID := mux.Vars(req)["id"]
	videoPath := fmt.Sprintf("%s%c%s.mp4", l.cfg.VideosPath, os.PathSeparator, videoID)
	http.ServeFile(w, req, videoPath)
}

func (l *StreamingApi) poster(w http.ResponseWriter, req *http.Request) {
	videoID := mux.Vars(req)["id"]
	posterPath := fmt.Sprintf("%s%c%s.jpg", l.cfg.PostersPath, os.PathSeparator, videoID)
	http.ServeFile(w, req, posterPath)
}

func writeErr(w http.ResponseWriter, err error, httpCode int) {
	logrus.Error(err.Error())

	// write error to response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var errMap = map[string]interface{}{
		"httpStatus": httpCode,
		"error":      err.Error(),
	}

	errJson, _ := json.Marshal(errMap)
	http.Error(w, string(errJson), httpCode)
}
