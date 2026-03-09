package service

import "github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/store"

type AnomalyService struct {
	store store.Repository
}

type DetectionService struct{}

func NewAnomalyService(repo store.Repository) *AnomalyService {
	return &AnomalyService{store: repo}
}

func NewDetectionService() *DetectionService {
	return &DetectionService{}
}

func (s *AnomalyService) List() []store.Anomaly {
	return s.store.ListAnomalies()
}

func (s *AnomalyService) ListBillingRecords() []store.BillingRecord {
	return s.store.ListBillingRecords()
}
