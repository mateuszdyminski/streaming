package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

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

// callStats holds various stats associated with HTTP request-response pair.
type callStats struct {
	w       http.ResponseWriter
	code    int // response status code
	resSize int64
}

func newCallStats(w http.ResponseWriter) *callStats {
	return &callStats{w: w}
}

func (r *callStats) Header() http.Header {
	return r.w.Header()
}

func (r *callStats) WriteHeader(code int) {
	r.w.WriteHeader(code)
	r.code = code
}

func (r *callStats) Write(p []byte) (n int, err error) {
	if r.code == 0 {
		r.code = http.StatusOK
	}
	n, err = r.w.Write(p)
	r.resSize += int64(n)
	return
}

func (r *callStats) StatusCode() int {
	return r.code
}

func (r *callStats) ResponseSize() int64 {
	return r.resSize
}

// NewLogginHandler creates and returns LoggingHandler with custom metrics.
func NewLogginHandler(root http.Handler) http.Handler {
	return LoggingHandler{root}
}

// LoggingHandler writes basic information about each request and response to
// the log.
type LoggingHandler struct {
	http.Handler
}

func (h LoggingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	stats := newCallStats(w)
	h.Handler.ServeHTTP(stats, req)

	elapsed := time.Since(t)

	buf := new(bytes.Buffer)
	buf.WriteString(getRemoteAddr(req))
	buf.WriteString(" - \"")
	buf.WriteString(req.Method)
	buf.WriteByte(' ')
	buf.WriteString(req.RequestURI)
	buf.WriteByte(' ')
	buf.WriteString(req.Proto)
	buf.WriteString("\" ")
	buf.WriteString(strconv.Itoa(stats.StatusCode()))
	buf.WriteByte(' ')
	buf.WriteString(strconv.FormatInt(stats.ResponseSize(), 10))
	buf.WriteString(" \"")
	buf.WriteString(req.UserAgent())
	buf.WriteString("\" Took: ")
	buf.WriteString(elapsed.String())

	logrus.Infof(buf.String())
}

// OptionsHandler allows all request with 'Options' http method to finish with status OK
type OptionsHandler struct {
	Handler http.Handler
}

func (h OptionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if http.MethodOptions == r.Method {
		w.WriteHeader(http.StatusOK)
		return
	}
	h.Handler.ServeHTTP(w, r)
}

func getRemoteAddr(r *http.Request) string {
	forwaredFor := r.Header.Get("X-Forwarded-For")
	if forwaredFor == "" {
		return r.RemoteAddr
	}

	return forwaredFor
}

func NewMetricsHandler(handler http.Handler) http.Handler {
	m := &metricsHandler{
		duration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "acme",
				Subsystem: "streaming_app",
				Name:      "http_durations_histogram_seconds",
				Help:      "Request time duration.",
			},
			[]string{"code", "method"},
		),
		requests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "acme",
				Subsystem: "streaming_app",
				Name:      "http_requests_total",
				Help:      "Total number of requests received.",
			},
			[]string{"code", "method"},
		),
	}

	prometheus.DefaultRegisterer.Register(m)

	return promhttp.InstrumentHandlerDuration(
		m.duration,
		promhttp.InstrumentHandlerCounter(
			m.requests,
			handler,
		),
	)
}

type metricsHandler struct {
	http.Handler
	duration *prometheus.HistogramVec
	requests *prometheus.CounterVec
}

// Describe implements prometheus.Collector interface.
func (d *metricsHandler) Describe(in chan<- *prometheus.Desc) {
	d.duration.Describe(in)
	d.requests.Describe(in)
}

// Collect implements prometheus.Collector interface.
func (d *metricsHandler) Collect(in chan<- prometheus.Metric) {
	d.duration.Collect(in)
	d.requests.Collect(in)
}
