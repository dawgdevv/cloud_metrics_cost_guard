package service

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/store"
)

func (s *DetectionService) Detect(job store.IngestionJob, records []store.BillingRecord) []store.Anomaly {
	buckets := make(map[string][]store.BillingRecord)
	for _, record := range records {
		buckets[record.Service] = append(buckets[record.Service], record)
	}

	anomalies := make([]store.Anomaly, 0)
	for serviceName, serviceRecords := range buckets {
		if len(serviceRecords) < 3 {
			continue
		}

		latest := serviceRecords[len(serviceRecords)-1]
		history := serviceRecords[:len(serviceRecords)-1]
		mean := average(history)
		stdDev := standardDeviation(history, mean)
		if stdDev == 0 {
			stdDev = 1
		}

		score := (latest.Amount - mean) / stdDev
		if score < 2 {
			continue
		}

		severity := "medium"
		if score >= 3 {
			severity = "high"
		}

		anomalies = append(anomalies, store.Anomaly{
			ID:            fmt.Sprintf("anomaly-%s-%d", sanitizeAnomalyID(serviceName), latest.UsageDate.Unix()),
			AccountID:     job.AccountID,
			Source:        job.Source,
			Service:       serviceName,
			CurrentSpend:  latest.Amount,
			ExpectedSpend: math.Round(mean*100) / 100,
			Score:         math.Round(score*100) / 100,
			Severity:      severity,
			DetectedAt:    time.Now().UTC(),
		})
	}

	return anomalies
}

func average(records []store.BillingRecord) float64 {
	total := 0.0
	for _, record := range records {
		total += record.Amount
	}
	return total / float64(len(records))
}

func standardDeviation(records []store.BillingRecord, mean float64) float64 {
	variance := 0.0
	for _, record := range records {
		delta := record.Amount - mean
		variance += delta * delta
	}
	return math.Sqrt(variance / float64(len(records)))
}

func sanitizeAnomalyID(value string) string {
	return strings.NewReplacer(" ", "-", "/", "-", "_", "-").Replace(strings.ToLower(value))
}
