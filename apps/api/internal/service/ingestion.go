package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/ingestion"
	appmetrics "github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/metrics"
	"github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/store"
)

type IngestionService struct {
	store    store.Repository
	detector *DetectionService
	metrics  *appmetrics.Collector
	sources  map[string]ingestion.Source
}

type QueueJobInput struct {
	AccountID string
	Days      int
	Source    string
	Scenario  string
}

func NewIngestionService(repo store.Repository, detector *DetectionService, metrics *appmetrics.Collector, sources map[string]ingestion.Source) *IngestionService {
	return &IngestionService{
		store:    repo,
		detector: detector,
		metrics:  metrics,
		sources:  sources,
	}
}

func (s *IngestionService) QueueJob(ctx context.Context, input QueueJobInput) (store.IngestionJob, error) {
	if input.Days <= 0 {
		input.Days = 7
	}
	input.Source = normalizeSource(input.Source)
	input.Scenario = normalizeScenario(input.Scenario)

	job := store.IngestionJob{
		ID:        fmt.Sprintf("job-%d", time.Now().UnixNano()),
		AccountID: input.AccountID,
		Days:      input.Days,
		Source:    input.Source,
		Scenario:  normalizeScenario(input.Scenario),
		Status:    "running",
		CreatedAt: time.Now().UTC(),
	}

	s.store.SaveJob(job)
	startedAt := time.Now()

	source, ok := s.sources[input.Source]
	if !ok {
		return s.failJob(job, startedAt, fmt.Errorf("unsupported source %q", input.Source))
	}

	records, err := source.Fetch(ctx, ingestion.FetchInput{
		AccountID: input.AccountID,
		Days:      input.Days,
		Scenario:  input.Scenario,
	})
	if err != nil {
		return s.failJob(job, startedAt, err)
	}

	s.store.SaveBillingRecords(records)
	s.metrics.ObserveBillingRecords(input.Source, len(records))

	severities := make([]string, 0)
	for _, anomaly := range s.detector.Detect(job, records) {
		s.store.SaveAnomaly(anomaly)
		severities = append(severities, anomaly.Severity)
	}
	s.metrics.ObserveAnomalies(input.Source, severities)

	completedAt := time.Now().UTC()
	job.Status = "completed"
	job.CompletedAt = &completedAt
	s.store.SaveJob(job)
	s.metrics.ObserveIngestion(input.Source, "completed", time.Since(startedAt))

	return job, nil
}

func (s *IngestionService) ListJobs() []store.IngestionJob {
	return s.store.ListJobs()
}

func (s *IngestionService) failJob(job store.IngestionJob, startedAt time.Time, err error) (store.IngestionJob, error) {
	completedAt := time.Now().UTC()
	job.Status = "failed"
	job.CompletedAt = &completedAt
	s.store.SaveJob(job)
	s.metrics.ObserveSourceFailure(job.Source)
	s.metrics.ObserveIngestion(job.Source, "failed", time.Since(startedAt))
	return job, err
}

func (s *IngestionService) ListBillingRecords() []store.BillingRecord {
	return s.store.ListBillingRecords()
}

func normalizeSource(source string) string {
	normalized := strings.TrimSpace(strings.ToLower(source))
	if normalized == "" {
		return "synthetic"
	}
	return normalized
}

func normalizeScenario(scenario string) string {
	normalized := strings.TrimSpace(strings.ToLower(scenario))
	switch normalized {
	case "gpu_spike", "storage_growth", "network_burst", "normal":
		return normalized
	default:
		return "gpu_spike"
	}
}
