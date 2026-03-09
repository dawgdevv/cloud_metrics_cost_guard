export function StatCards({ anomalies, jobs, records }) {
  return (
    <section className="stats">
      <article className="stat-card">
        <span className="stat-label">Accounts Monitored</span>
        <strong className="stat-value">1</strong>
      </article>
      <article className="stat-card">
        <span className="stat-label">Open Anomalies</span>
        <strong className="stat-value accent">{anomalies.length}</strong>
      </article>
      <article className="stat-card">
        <span className="stat-label">Jobs Created</span>
        <strong className="stat-value">{jobs.length}</strong>
      </article>
      <article className="stat-card">
        <span className="stat-label">Synthetic Records</span>
        <strong className="stat-value">{records.length}</strong>
      </article>
    </section>
  );
}