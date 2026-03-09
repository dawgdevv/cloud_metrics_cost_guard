package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/ingestion"
	"github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/service"
	"github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/store"
)

type Handler struct {
	deps      RouterDependencies
	tokenAuth *jwtauth.JWTAuth
}

type ingestRequest struct {
	AccountID string `json:"account_id"`
	Days      int    `json:"days"`
	Source    string `json:"source"`
	Scenario  string `json:"scenario"`
}

func NewHandler(deps RouterDependencies, tokenAuth *jwtauth.JWTAuth) *Handler {
	return &Handler{deps: deps, tokenAuth: tokenAuth}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) MetricsSummary(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.deps.Metrics.Snapshot())
}

func (h *Handler) IssueDemoToken(w http.ResponseWriter, _ *http.Request) {
	_, tokenString, err := h.tokenAuth.Encode(map[string]any{
		"sub":  "demo-user",
		"role": "operator",
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not issue token"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

func (h *Handler) CreateIngestionJob(w http.ResponseWriter, r *http.Request) {
	var request ingestRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	job, err := h.deps.IngestionService.QueueJob(r.Context(), service.QueueJobInput{
		AccountID: request.AccountID,
		Days:      request.Days,
		Source:    request.Source,
		Scenario:  request.Scenario,
	})
	if err != nil {
		statusCode := http.StatusBadRequest
		payload := map[string]any{"error": map[string]string{"message": err.Error()}, "job": job}

		var sourceErr *ingestion.SourceError
		if errors.As(err, &sourceErr) {
			statusCode = http.StatusBadGateway
			payload["error"] = sourceErr.PublicError()
		}

		writeJSON(w, statusCode, payload)
		return
	}

	writeJSON(w, http.StatusAccepted, job)
}

func (h *Handler) ListAnomalies(w http.ResponseWriter, _ *http.Request) {
	anomalies := h.deps.AnomalyService.List()
	writeJSON(w, http.StatusOK, map[string][]store.Anomaly{"items": anomalies})
}

func (h *Handler) ListJobs(w http.ResponseWriter, _ *http.Request) {
	jobs := h.deps.IngestionService.ListJobs()
	writeJSON(w, http.StatusOK, map[string][]store.IngestionJob{"items": jobs})
}

func (h *Handler) ListBillingRecords(w http.ResponseWriter, _ *http.Request) {
	records := h.deps.AnomalyService.ListBillingRecords()
	writeJSON(w, http.StatusOK, map[string][]store.BillingRecord{"items": records})
}
