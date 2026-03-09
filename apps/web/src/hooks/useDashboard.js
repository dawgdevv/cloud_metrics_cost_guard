import { useEffect, useState } from "react";
import { INITIAL_INGESTION_FORM } from "../constants/ingestion";
import { ApiRequestError } from "../services/apiClient";
import {
  createIngestionJob,
  getDemoToken,
  getMetricsSummary,
  listAnomalies,
  listBillingRecords,
  listJobs,
} from "../services/dashboardApi";

export function useDashboard() {
  const [token, setToken] = useState("");
  const [form, setForm] = useState(INITIAL_INGESTION_FORM);
  const [activeSource, setActiveSource] = useState(INITIAL_INGESTION_FORM.source);
  const [jobs, setJobs] = useState([]);
  const [anomalies, setAnomalies] = useState([]);
  const [records, setRecords] = useState([]);
  const [metrics, setMetrics] = useState(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [ingestionError, setIngestionError] = useState(null);

  async function refreshDashboard(nextToken) {
    const [jobResponse, anomalyResponse, recordResponse, metricsResponse] = await Promise.all([
      listJobs(nextToken),
      listAnomalies(nextToken),
      listBillingRecords(nextToken),
      getMetricsSummary(nextToken),
    ]);

    setJobs(jobResponse.items || []);
    setAnomalies(anomalyResponse.items || []);
    setRecords(recordResponse.items || []);
    setMetrics(metricsResponse || null);
  }

  useEffect(() => {
    async function bootstrap() {
      try {
        const tokenResponse = await getDemoToken();
        const nextToken = tokenResponse.token;
        setToken(nextToken);
        await refreshDashboard(nextToken);
        setError("");
        setIngestionError(null);
      } catch (err) {
        setError(err.message || "Failed to bootstrap dashboard");
      } finally {
        setLoading(false);
      }
    }

    bootstrap();
  }, []);

  async function submitIngestion() {
    if (!token || submitting) {
      return;
    }

    setSubmitting(true);
    setIngestionError(null);
    try {
      await createIngestionJob(form, token);
      await refreshDashboard(token);
      setError("");
      setIngestionError(null);
    } catch (err) {
      if (err instanceof ApiRequestError) {
        const nextIngestionError = {
          source: form.source,
          status: err.status,
          code: err.code,
          message: err.message || "Failed to create job",
          hint: err.hint || "",
          details: err.details || "",
          job: err.job || null,
        };
        setIngestionError(nextIngestionError);
        setError(form.source === "aws" ? "" : nextIngestionError.message);
      } else {
        const fallbackMessage = err.message || "Failed to create job";
        setIngestionError({
          source: form.source,
          status: 0,
          code: "request_failed",
          message: fallbackMessage,
          hint: "",
          details: "",
          job: null,
        });
        setError(form.source === "aws" ? "" : fallbackMessage);
      }
    } finally {
      setSubmitting(false);
    }
  }

  return {
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
  };
}