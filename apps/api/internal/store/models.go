package store

import "time"

type IngestionJob struct {
	ID          string     `json:"id"`
	AccountID   string     `json:"account_id"`
	Days        int        `json:"days"`
	Source      string     `json:"source"`
	Scenario    string     `json:"scenario"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type Anomaly struct {
	ID            string    `json:"id"`
	AccountID     string    `json:"account_id"`
	Source        string    `json:"source"`
	Service       string    `json:"service"`
	CurrentSpend  float64   `json:"current_spend"`
	ExpectedSpend float64   `json:"expected_spend"`
	Score         float64   `json:"score"`
	Severity      string    `json:"severity"`
	DetectedAt    time.Time `json:"detected_at"`
}

type BillingRecord struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Service   string    `json:"service"`
	UsageDate time.Time `json:"usage_date"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Source    string    `json:"source"`
	Scenario  string    `json:"scenario"`
}

type Repository interface {
	SaveJob(job IngestionJob)
	ListJobs() []IngestionJob
	SaveAnomaly(anomaly Anomaly)
	ListAnomalies() []Anomaly
	SaveBillingRecords(records []BillingRecord)
	ListBillingRecords() []BillingRecord
}
