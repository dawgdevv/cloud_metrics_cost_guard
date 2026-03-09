package ingestion

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/store"
)

type SyntheticSource struct{}

func NewSyntheticSource() *SyntheticSource {
	return &SyntheticSource{}
}

func (s *SyntheticSource) Fetch(_ context.Context, input FetchInput) ([]store.BillingRecord, error) {
	if input.Days <= 0 {
		input.Days = 7
	}

	services := []string{"Amazon EC2 GPU", "Amazon S3", "AWS Data Transfer"}
	records := make([]store.BillingRecord, 0, input.Days*len(services))
	dayZero := time.Now().UTC()
	scenario := normalizeScenario(input.Scenario)

	for offset := input.Days - 1; offset >= 0; offset-- {
		usageDate := time.Date(dayZero.Year(), dayZero.Month(), dayZero.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -offset)
		isLatestDay := offset == 0

		for serviceIndex, serviceName := range services {
			amount := baselineSpend(serviceName, serviceIndex, offset)
			amount = applyScenario(scenario, serviceName, amount, isLatestDay)

			records = append(records, store.BillingRecord{
				ID:        fmt.Sprintf("record-%s-%d-%d", sanitizeID(input.AccountID), usageDate.Unix(), serviceIndex),
				AccountID: input.AccountID,
				Service:   serviceName,
				UsageDate: usageDate,
				Amount:    math.Round(amount*100) / 100,
				Currency:  "USD",
				Source:    "synthetic",
				Scenario:  scenario,
			})
		}
	}

	return records, nil
}

func baselineSpend(serviceName string, serviceIndex int, offset int) float64 {
	base := map[string]float64{
		"Amazon EC2 GPU":    520,
		"Amazon S3":         88,
		"AWS Data Transfer": 46,
	}

	return base[serviceName] + float64((serviceIndex+1)*offset%5)*6 + float64(offset%3)*4
}

func applyScenario(scenario string, serviceName string, amount float64, isLatestDay bool) float64 {
	if !isLatestDay {
		return amount
	}

	switch scenario {
	case "gpu_spike":
		if serviceName == "Amazon EC2 GPU" {
			return amount * 1.82
		}
	case "storage_growth":
		if serviceName == "Amazon S3" {
			return amount * 2.35
		}
	case "network_burst":
		if serviceName == "AWS Data Transfer" {
			return amount * 2.9
		}
	case "normal":
		return amount * 1.03
	}

	return amount
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

func sanitizeID(value string) string {
	return strings.NewReplacer(" ", "-", "/", "-", "_", "-").Replace(strings.ToLower(value))
}
