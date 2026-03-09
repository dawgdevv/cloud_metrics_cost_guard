package metrics

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Collector struct {
	mu                     sync.RWMutex
	registry               *prometheus.Registry
	httpRequestsTotal      *prometheus.CounterVec
	ingestionJobsTotal     *prometheus.CounterVec
	billingRecordsTotal    *prometheus.CounterVec
	anomaliesDetectedTotal *prometheus.CounterVec
	ingestionDuration      *prometheus.HistogramVec
	sourceFailuresTotal    *prometheus.CounterVec
	appUp                  prometheus.Gauge
	summary                Summary
}

type Summary struct {
	AppUp                bool               `json:"app_up"`
	HTTPRequestsTotal    int                `json:"http_requests_total"`
	SyntheticJobsTotal   int                `json:"synthetic_jobs_total"`
	BillingRecordsTotal  int                `json:"billing_records_total"`
	AnomaliesDetected    int                `json:"anomalies_detected_total"`
	AnomaliesBySeverity  map[string]int     `json:"anomalies_by_severity"`
	SourceFailures       map[string]int     `json:"source_failures"`
	LastIngestionSeconds map[string]float64 `json:"last_ingestion_duration_seconds"`
}

func NewCollector() *Collector {
	registry := prometheus.NewRegistry()
	collector := &Collector{
		registry: registry,
		httpRequestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests handled by the API.",
		}, []string{"method", "path", "status"}),
		ingestionJobsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "ingestion_jobs_total",
			Help: "Total ingestion jobs by source and final status.",
		}, []string{"source", "status"}),
		billingRecordsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "billing_records_ingested_total",
			Help: "Total billing records ingested by source.",
		}, []string{"source"}),
		anomaliesDetectedTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "anomalies_detected_total",
			Help: "Total anomalies detected by source and severity.",
		}, []string{"source", "severity"}),
		ingestionDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "ingestion_duration_seconds",
			Help:    "Duration of ingestion jobs by source.",
			Buckets: prometheus.DefBuckets,
		}, []string{"source", "status"}),
		sourceFailuresTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "source_fetch_failures_total",
			Help: "Total source fetch failures by source.",
		}, []string{"source"}),
		appUp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "app_up",
			Help: "Application liveness gauge.",
		}),
	}

	registry.MustRegister(
		collector.httpRequestsTotal,
		collector.ingestionJobsTotal,
		collector.billingRecordsTotal,
		collector.anomaliesDetectedTotal,
		collector.ingestionDuration,
		collector.sourceFailuresTotal,
		collector.appUp,
	)
	collector.appUp.Set(1)
	collector.summary = Summary{
		AppUp:                true,
		AnomaliesBySeverity:  map[string]int{"high": 0, "medium": 0},
		SourceFailures:       map[string]int{"synthetic": 0, "aws": 0},
		LastIngestionSeconds: map[string]float64{"synthetic": 0, "aws": 0},
	}

	return collector
}

func (c *Collector) Handler() http.Handler {
	return promhttp.HandlerFor(c.registry, promhttp.HandlerOpts{})
}

func (c *Collector) ObserveHTTPRequest(method string, path string, statusCode int) {
	c.httpRequestsTotal.WithLabelValues(method, path, strconv.Itoa(statusCode)).Inc()
	c.mu.Lock()
	c.summary.HTTPRequestsTotal++
	c.mu.Unlock()
}

func (c *Collector) ObserveIngestion(source string, status string, duration time.Duration) {
	c.ingestionJobsTotal.WithLabelValues(source, status).Inc()
	c.ingestionDuration.WithLabelValues(source, status).Observe(duration.Seconds())
	c.mu.Lock()
	if source == "synthetic" && status == "completed" {
		c.summary.SyntheticJobsTotal++
	}
	c.summary.LastIngestionSeconds[source] = duration.Seconds()
	c.mu.Unlock()
}

func (c *Collector) ObserveBillingRecords(source string, count int) {
	c.billingRecordsTotal.WithLabelValues(source).Add(float64(count))
	c.mu.Lock()
	if source == "synthetic" {
		c.summary.BillingRecordsTotal += count
	}
	c.mu.Unlock()
}

func (c *Collector) ObserveAnomalies(source string, anomalies []string) {
	for _, severity := range anomalies {
		c.anomaliesDetectedTotal.WithLabelValues(source, severity).Inc()
	}
	if len(anomalies) == 0 {
		return
	}
	c.anomaliesDetectedTotal.WithLabelValues(source, "all").Add(float64(len(anomalies)))
	c.mu.Lock()
	c.summary.AnomaliesDetected += len(anomalies)
	for _, severity := range anomalies {
		c.summary.AnomaliesBySeverity[severity]++
	}
	c.mu.Unlock()
}

func (c *Collector) ObserveSourceFailure(source string) {
	c.sourceFailuresTotal.WithLabelValues(source).Inc()
	c.mu.Lock()
	c.summary.SourceFailures[source]++
	c.mu.Unlock()
}

func (c *Collector) Snapshot() Summary {
	c.mu.RLock()
	defer c.mu.RUnlock()

	summary := c.summary
	summary.AnomaliesBySeverity = cloneIntMap(c.summary.AnomaliesBySeverity)
	summary.SourceFailures = cloneIntMap(c.summary.SourceFailures)
	summary.LastIngestionSeconds = cloneFloatMap(c.summary.LastIngestionSeconds)
	return summary
}

func cloneIntMap(input map[string]int) map[string]int {
	output := make(map[string]int, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}

func cloneFloatMap(input map[string]float64) map[string]float64 {
	output := make(map[string]float64, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}
