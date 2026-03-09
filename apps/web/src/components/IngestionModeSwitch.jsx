import { INGESTION_SOURCE_STATUS } from "../constants/ingestion";

export function IngestionModeSwitch({ activeSource, onChange }) {
  return (
    <section className="mode-switch" aria-label="Ingestion mode switch">
      {Object.entries(INGESTION_SOURCE_STATUS).map(([source, config]) => (
        <button
          key={source}
          type="button"
          className={`mode-button ${activeSource === source ? "active" : ""}`}
          onClick={() => onChange(source)}
        >
          <span>{config.label}</span>
          <small>{config.available ? "Available" : "Stub"}</small>
        </button>
      ))}
    </section>
  );
}