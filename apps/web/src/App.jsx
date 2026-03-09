import { AnomalyTable } from "./components/AnomalyTable";
import { AWSIngestionPanel } from "./components/AWSIngestionPanel";
import { BillingRecordsTable } from "./components/BillingRecordsTable";
import { DemoIngestionForm } from "./components/DemoIngestionForm";
import { IngestionModeSwitch } from "./components/IngestionModeSwitch";
import { JobList } from "./components/JobList";
import { PrometheusStatusCard } from "./components/PrometheusStatusCard";
import { StatCards } from "./components/StatCards";
import { useDashboard } from "./hooks/useDashboard";

export default function App() {
  const {
    form,
    setForm,
    activeSource,
    setActiveSource,
    jobs,
    anomalies,
    records,
    metrics,
    loading,
    submitting,
    error,
    ingestionError,
    submitIngestion,
  } = useDashboard();

  return (
    <div className="shell">
      <header className="hero">
        <p className="eyebrow">Cloud FinOps Platform</p>
        <h1><span>Cost Guard</span></h1>
        <p className="lede">
          Demo mode is active. Generate synthetic billing data, run ingestion locally, and detect
          anomalies without AWS billing access.
        </p>
      </header>

      <StatCards anomalies={anomalies} jobs={jobs} records={records} />

      {error ? <p className="card-copy">{error}</p> : null}

      <IngestionModeSwitch
        activeSource={activeSource}
        onChange={(source) => {
          setActiveSource(source);
          setForm((current) => ({ ...current, source }));
        }}
      />

      <section className="layout-grid">
        {activeSource === "synthetic" ? (
          <DemoIngestionForm
            form={form}
            setForm={setForm}
            submitting={submitting}
            onSubmit={submitIngestion}
          />
        ) : (
          <AWSIngestionPanel
            form={form}
            setForm={setForm}
            submitting={submitting}
            onSubmit={submitIngestion}
            error={activeSource === "aws" ? ingestionError : null}
          />
        )}

        <article className="card">
          <div className="card-header">
            <h2>Recent Jobs</h2>
          </div>
          <JobList jobs={jobs} loading={loading} />
        </article>
      </section>

      <section className="layout-grid layout-grid-triple">
        <PrometheusStatusCard metrics={metrics} loading={loading} />
        <article className="card">
          <div className="card-header">
            <h2>Source State</h2>
          </div>
          <p className="card-copy">
            Both ingestion paths are wired. Demo mode generates synthetic records locally, while AWS mode
            reads live daily spend from Cost Explorer using the credentials available to the API.
          </p>
        </article>
      </section>

      <section className="card anomalies-section">
        <div className="card-header">
          <h2>Detected Anomalies</h2>
        </div>
        <AnomalyTable anomalies={anomalies} loading={loading} />
      </section>

      <section className="card anomalies-section">
        <div className="card-header">
          <h2>Generated Billing Records</h2>
        </div>
        <BillingRecordsTable records={records} loading={loading} />
      </section>
    </div>
  );
}
