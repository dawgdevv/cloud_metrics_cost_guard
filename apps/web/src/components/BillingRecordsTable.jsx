import { SCENARIO_LABELS } from "../constants/ingestion";
import { LoadingState } from "./LoadingState";

export function BillingRecordsTable({ records, loading }) {
  if (loading) {
    return <LoadingState />;
  }

  return (
    <div className="table-wrapper">
      <table>
        <thead>
          <tr>
            <th>Date</th>
            <th>Service</th>
            <th>Amount</th>
            <th>Source</th>
            <th>Scenario</th>
          </tr>
        </thead>
        <tbody>
          {records.slice(0, 12).map((record) => (
            <tr key={record.id}>
              <td>{new Date(record.usage_date).toLocaleDateString()}</td>
              <td>{record.service}</td>
              <td>${record.amount?.toFixed(2)}</td>
              <td>{record.source}</td>
              <td>{SCENARIO_LABELS[record.scenario] || record.scenario}</td>
            </tr>
          ))}
          {!loading && records.length === 0 ? (
            <tr>
              <td colSpan="5">No billing records yet. Run a synthetic ingestion job.</td>
            </tr>
          ) : null}
        </tbody>
      </table>
    </div>
  );
}