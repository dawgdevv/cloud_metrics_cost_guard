function formatSeconds(value) {
  if (!value) {
    return "0.00s";
  }
  return `${value.toFixed(2)}s`;
}

export function PrometheusStatusCard({ metrics, loading }) {
  return (
    <article className="card">
      <div className="card-header">
        <h2>Prometheus Status</h2>
        <span className={`panel-badge ${metrics?.app_up ? "panel-badge-live" : ""}`}>
          {metrics?.app_up ? "Live" : loading ? "Loading" : "Unknown"}
        </span>
      </div>
      <p className="card-copy">
        Live backend snapshot of key Prometheus counters for quick operator visibility.
      </p>

      <div className="metrics-grid">
        <div className="metric-row">
          <span>HTTP Requests</span>
          <strong>{metrics?.http_requests_total ?? 0}</strong>
        </div>
        <div className="metric-row">
          <span>Synthetic Jobs</span>
          <strong>{metrics?.synthetic_jobs_total ?? 0}</strong>
        </div>
        <div className="metric-row">
          <span>Records Ingested</span>
          <strong>{metrics?.billing_records_total ?? 0}</strong>
        </div>
        <div className="metric-row">
          <span>Anomalies Detected</span>
          <strong>{metrics?.anomalies_detected_total ?? 0}</strong>
        </div>
        <div className="metric-row">
          <span>High Severity</span>
          <strong>{metrics?.anomalies_by_severity?.high ?? 0}</strong>
        </div>
        <div className="metric-row">
          <span>Last Ingestion</span>
          <strong>{formatSeconds(metrics?.last_ingestion_duration_seconds?.synthetic)}</strong>
        </div>
      </div>
    </article>
  );
}