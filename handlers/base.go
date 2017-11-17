package handlers

import (
	"log"
	"net/http"
	"time"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

// LogRequest is a middleware that logs details about the request/response.
func LogRequest() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			t := time.Now()
			next.ServeHTTP(res, req)
			log.Printf("%s \"%s %s %s\" \"%s\" \"%s\" \"Took: %s\"\n", req.RemoteAddr,
				req.Method, req.RequestURI, req.Proto, req.Referer(), req.UserAgent(), time.Since(t))
		})
	}
}

// WriteJSON write response to client, response is a struct defining JSON reply.
func WriteJSON(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json, err := json.Marshal(response)
	if err != nil {
		return err
	}

	if _, err := w.Write(json); err != nil {
		return err
	}

	return nil
}

func WriteErr(w http.ResponseWriter, err error, httpCode int) {
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
