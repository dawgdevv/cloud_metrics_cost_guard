CREATE TABLE IF NOT EXISTS ingestion_jobs (
    id TEXT PRIMARY KEY,
    account_id TEXT NOT NULL,
    days INTEGER NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS billing_records (
    id BIGSERIAL PRIMARY KEY,
    account_id TEXT NOT NULL,
    service TEXT NOT NULL,
    usage_date DATE NOT NULL,
    amount NUMERIC(12, 2) NOT NULL,
    currency TEXT NOT NULL DEFAULT 'USD',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS anomalies (
    id TEXT PRIMARY KEY,
    account_id TEXT NOT NULL,
    service TEXT NOT NULL,
    current_spend NUMERIC(12, 2) NOT NULL,
    expected_spend NUMERIC(12, 2) NOT NULL,
    score NUMERIC(8, 2) NOT NULL,
    severity TEXT NOT NULL,
    detected_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_billing_records_account_date
    ON billing_records (account_id, usage_date DESC);

CREATE INDEX IF NOT EXISTS idx_anomalies_detected_at
    ON anomalies (detected_at DESC);
