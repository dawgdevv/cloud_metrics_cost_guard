import { INGESTION_SOURCE_STATUS } from "../constants/ingestion";

export function AWSIngestionPanel({ form, setForm, submitting, onSubmit, error }) {
  const awsState = INGESTION_SOURCE_STATUS.aws;

  function handleChange(field, transform = (value) => value) {
    return (event) => {
      setForm((current) => ({
        ...current,
        source: "aws",
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
        <h2>AWS Ingestion</h2>
        <span className="panel-badge">{awsState.available ? "Ready" : "Stub"}</span>
      </div>
      <p className="card-copy">
        {awsState.description}
      </p>

      {error ? (
        <section className="aws-error-panel" aria-live="polite">
          <div className="aws-error-header">
            <strong>AWS ingestion failed</strong>
            {error.code ? <span className="aws-error-code">{error.code}</span> : null}
          </div>
          <p className="aws-error-message">{error.message}</p>
          {error.hint ? <p className="aws-error-hint">{error.hint}</p> : null}
          {error.details ? (
            <details className="aws-error-details">
              <summary>Raw AWS error details</summary>
              <pre>{error.details}</pre>
            </details>
          ) : null}
        </section>
      ) : (
        <section className="aws-help-panel">
          <strong>Before you run it</strong>
          <p>
            The API process needs valid AWS credentials and permission to call Cost Explorer.
            If this fails, the panel will show whether the problem is credentials, IAM access,
            region configuration, or request validation.
          </p>
        </section>
      )}

      <form className="form" onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="aws_account_id">AWS Account ID</label>
          <input
            id="aws_account_id"
            type="text"
            value={form.account_id}
            onChange={handleChange("account_id")}
            placeholder="123456789012"
          />
        </div>
        <div className="form-group">
          <label htmlFor="aws_days">Lookback Days</label>
          <input
            id="aws_days"
            type="number"
            min="1"
            max="90"
            value={form.days}
            onChange={handleChange("days", Number)}
          />
        </div>
        <div className="form-group">
          <label htmlFor="aws_source">Source</label>
          <input id="aws_source" type="text" value="aws" disabled readOnly />
        </div>
        <button className="btn" type="submit" disabled={submitting}>
          {submitting ? "Creating..." : "Run AWS Ingestion"}
        </button>
      </form>
    </article>
  );
}