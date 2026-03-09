import { SCENARIO_LABELS } from "../constants/ingestion";

export function DemoIngestionForm({ form, setForm, submitting, onSubmit }) {
  function handleChange(field, transform = (value) => value) {
    return (event) => {
      setForm((current) => ({
        ...current,
        source: "synthetic",
        [field]: transform(event.target.value),
      }));
    };
  }

  function handleSubmit(event) {
    event.preventDefault();
    onSubmit();
  }

  return (
    <article className="card">
      <div className="card-header">
        <h2>Demo Ingestion</h2>
      </div>
      <p className="card-copy">
        Synthetic-only workflow for testing and demos. This panel generates sample billing series and
        runs anomaly detection locally.
      </p>
      <form className="form" onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="demo_account_id">Account ID</label>
          <input
            id="demo_account_id"
            type="text"
            value={form.account_id}
            onChange={handleChange("account_id")}
            placeholder="123456789012"
          />
        </div>

        <div className="form-group">
          <label htmlFor="demo_scenario">Scenario</label>
          <select id="demo_scenario" value={form.scenario} onChange={handleChange("scenario")}>
            {Object.entries(SCENARIO_LABELS).map(([value, label]) => (
              <option key={value} value={value}>
                {label}
              </option>
            ))}
          </select>
        </div>

        <div className="form-group">
          <label htmlFor="demo_days">Lookback Days</label>
          <input
            id="demo_days"
            type="number"
            min="1"
            max="90"
            value={form.days}
            onChange={handleChange("days", Number)}
          />
        </div>

        <button className="btn" type="submit" disabled={submitting}>
          {submitting ? "Creating..." : "Run Demo Ingestion"}
        </button>
      </form>
    </article>
  );
}