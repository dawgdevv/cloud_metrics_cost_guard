import { LoadingState } from "./LoadingState";

export function JobList({ jobs, loading }) {
  if (loading) {
    return <LoadingState />;
  }

  if (jobs.length === 0) {
    return <div className="empty-state">No ingestion jobs yet. Create one to get started.</div>;
  }

  return (
    <div className="job-list">
      {jobs.map((job) => (
        <div className="job-item" key={job.id}>
          <div className="job-info">
            <span className="job-account">{job.account_id}</span>
            <span className="job-id">{job.id.slice(0, 18)}...</span>
          </div>
          <span className={`job-status ${job.status}`}>{job.status}</span>
        </div>
      ))}
    </div>
  );
}