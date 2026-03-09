package ingestion

import (
	"context"

	"github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/store"
)

type FetchInput struct {
	AccountID string
	Days      int
	Scenario  string
}

type Source interface {
	Fetch(ctx context.Context, input FetchInput) ([]store.BillingRecord, error)
}
