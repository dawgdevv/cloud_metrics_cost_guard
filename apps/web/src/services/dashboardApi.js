import { apiRequest } from "./apiClient";

function authHeaders(token, extraHeaders = {}) {
  return {
    ...extraHeaders,
    Authorization: `Bearer ${token}`,
  };
}

export function getDemoToken() {
  return apiRequest("/auth/token", { method: "POST" });
}

export function createIngestionJob(payload, token) {
  return apiRequest("/api/v1/ingest", {
    method: "POST",
    headers: authHeaders(token, { "Content-Type": "application/json" }),
    body: JSON.stringify(payload),
  });
}

export function listJobs(token) {
  return apiRequest("/api/v1/jobs", {
    headers: authHeaders(token),
  });
}

export function listAnomalies(token) {
  return apiRequest("/api/v1/anomalies", {
    headers: authHeaders(token),
  });
}

export function listBillingRecords(token) {
  return apiRequest("/api/v1/billing-records", {
    headers: authHeaders(token),
  });
}

export function getMetricsSummary(token) {
  return apiRequest("/api/v1/metrics/summary", {
    headers: authHeaders(token),
  });
}