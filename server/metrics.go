package server

import (
	"net/http"
	"time"

	"github.com/etzba/pggo/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
)

var labels = []string{"endpoint", "method", "content_type"}

type Shipper interface {
	Register()
	Collect(start time.Time, r *http.Request)
	NewGauge(name, help string, labels []string) *prometheus.GaugeVec
	NewCounter(name, help string, labels []string) *prometheus.CounterVec
	NewHistogram(name, help string, labels []string) *prometheus.HistogramVec
}

// NewShipper creates new prometheus api client to push metrics
func NewShipper(logger *logger.Log) Shipper {
	return &collector{
		Logger:    logger,
		Namespace: "gopu",
	}
}

type collector struct {
	Logger    *logger.Log
	Namespace string
	Metrics   Metrics
}

func (c *collector) Register() {
	c.Logger.Info("register metrics")
	c.Metrics.handlerDuration = c.NewHistogram("request_duration", "measure the time of request until response in handler", labels)
	c.Metrics.httpRequestCount = c.NewCounter("total_requests", "count requests to endpoint", labels)

	prometheus.MustRegister(c.Metrics.handlerDuration)
	prometheus.MustRegister(c.Metrics.httpRequestCount)
}

// NewGauge https://prometheus.io/docs/concepts/metric_types/#gauge
// A gauge is a metric that represents a single numerical value that can arbitrarily go up and down.
func (c *collector) NewGauge(name, help string, labels []string) *prometheus.GaugeVec {
	gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: c.Namespace,
		Subsystem: "apiserver",
		Name:      name,
		Help:      help,
	}, labels)
	return gauge
}

// NewCounter https://prometheus.io/docs/concepts/metric_types/#counter
// A counter is a cumulative metric that represents a single monotonically increasing counter whose value can only increase or be reset to zero on restart.
// Use to count total executions from a result
func (c *collector) NewCounter(name, help string, labels []string) *prometheus.CounterVec {
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: c.Namespace,
		Subsystem: "apiserver",
		Name:      name,
		Help:      help,
	}, labels)
	return counter
}

// NewHistogram https://prometheus.io/docs/concepts/metric_types/#histogram
// A histogram samples observations (usually things like request durations or response sizes) and counts them in configurable buckets. It also provides a sum of all observed values.
// The best fit for this tool to collect results.
func (c *collector) NewHistogram(name, help string, labels []string) *prometheus.HistogramVec {
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: c.Namespace,
		Subsystem: "apiserver",
		Name:      name,
		Help:      help,
	}, labels)
	return histogram
}

type Metrics struct {
	httpRequestCount *prometheus.CounterVec
	handlerDuration  *prometheus.HistogramVec
}

func (c *collector) Collect(start time.Time, r *http.Request) {
	c.Metrics.httpRequestCount.WithLabelValues(r.URL.Path, r.Method, "Content-Type").Inc()
	c.Metrics.handlerDuration.WithLabelValues(r.URL.Path, r.Method, "Content-Type").Observe(float64(time.Since(start)))
}
