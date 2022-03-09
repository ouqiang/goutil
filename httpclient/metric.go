package httpclient

import (
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var histogramBuckets = []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 3, 5, 10}

// Metric 统计
type Metric struct {
	namespace                        string
	subsystem                        string
	formatUrl                        FormatUrl
	httpClientRequestTotal           *prometheus.CounterVec
	httpClientRequestDurationSeconds *prometheus.HistogramVec
}

type FormatUrl func(u *url.URL)

func NewMetric(namespace, subsystem string, formatUrl FormatUrl) *Metric {
	m := &Metric{
		namespace: namespace,
		subsystem: subsystem,
		formatUrl: formatUrl,
	}
	if m.formatUrl == nil {
		m.formatUrl = func(url *url.URL) {}
	}

	m.httpClientRequestDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "http_client_request_duration_seconds",
		Help:      "http client request duration seconds",
		Buckets:   histogramBuckets,
	}, []string{"host", "path"})
	prometheus.MustRegister(m.httpClientRequestDurationSeconds)

	m.httpClientRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "http_client_request_total",
		Help:      "http client request total",
	}, []string{"host", "path", "status"})
	prometheus.MustRegister(m.httpClientRequestTotal)

	return m
}

func (m *Metric) Count(url *url.URL, err error) {
	var status string
	if err != nil {
		status = "failure"
	} else {
		status = "success"
	}
	m.formatUrl(url)
	m.httpClientRequestTotal.WithLabelValues(
		url.Host,
		url.Path,
		status,
	).Inc()
}

func (m *Metric) Latency(url *url.URL, d time.Duration) {
	m.formatUrl(url)

	m.httpClientRequestDurationSeconds.WithLabelValues(url.Host, url.Path).Observe(d.Seconds())
}
