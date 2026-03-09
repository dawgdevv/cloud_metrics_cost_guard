import { LoadingState } from "./LoadingState";

export function AnomalyTable({ anomalies, loading }) {
  if (loading) {
    return <LoadingState />;
  }

  if (anomalies.length === 0) {
    return <div className="empty-state">No anomalies detected. Your cloud costs look healthy.</div>;
  }

  return (
    <div className="table-wrapper">
      <table>
        <thead>
          <tr>
            <th>Service</th>
            <th>Account</th>
            <th>Current</th>
            <th>Expected</th>
            <th>Score</th>
            <th>Severity</th>
          </tr>
        </thead>
        <tbody>
          {anomalies.map((anomaly) => (
            <tr key={anomaly.id}>
              <td>{anomaly.service}</td>
              <td>{anomaly.account_id?.slice(0, 12)}...</td>
              <td>${anomaly.current_spend?.toFixed(2)}</td>
              <td>${anomaly.expected_spend?.toFixed(2)}</td>
              <td>{anomaly.score?.toFixed(1)}</td>
              <td>
                <span className={`severity ${anomaly.severity?.toLowerCase()}`}>
                  {anomaly.severity}
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}