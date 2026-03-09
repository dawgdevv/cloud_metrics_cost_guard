package store

import "sync"

type MemoryStore struct {
	mu             sync.RWMutex
	jobs           []IngestionJob
	anomalies      []Anomaly
	billingRecords []BillingRecord
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		jobs:           make([]IngestionJob, 0),
		anomalies:      make([]Anomaly, 0),
		billingRecords: make([]BillingRecord, 0),
	}
}

func (s *MemoryStore) SaveJob(job IngestionJob) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for index := range s.jobs {
		if s.jobs[index].ID == job.ID {
			s.jobs[index] = job
			return
		}
	}
	s.jobs = append([]IngestionJob{job}, s.jobs...)
}

func (s *MemoryStore) ListJobs() []IngestionJob {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]IngestionJob, len(s.jobs))
	copy(items, s.jobs)
	return items
}

func (s *MemoryStore) SaveAnomaly(anomaly Anomaly) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.anomalies = append([]Anomaly{anomaly}, s.anomalies...)
}

func (s *MemoryStore) ListAnomalies() []Anomaly {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]Anomaly, len(s.anomalies))
	copy(items, s.anomalies)
	return items
}

func (s *MemoryStore) SaveBillingRecords(records []BillingRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.billingRecords = append(records, s.billingRecords...)
}

func (s *MemoryStore) ListBillingRecords() []BillingRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]BillingRecord, len(s.billingRecords))
	copy(items, s.billingRecords)
	return items
}
