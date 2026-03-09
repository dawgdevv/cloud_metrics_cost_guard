export function LoadingState({ label = "Initializing" }) {
  return (
    <div className="loading">
      <span className="loading-dot"></span>
      <span className="loading-dot"></span>
      <span className="loading-dot"></span>
      <span>{label}</span>
    </div>
  );
}