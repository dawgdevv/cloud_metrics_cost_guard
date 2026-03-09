package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/config"
	custommiddleware "github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/http/middleware"
	appmetrics "github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/metrics"
	"github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/service"
)

type RouterDependencies struct {
	Config           config.Config
	Metrics          *appmetrics.Collector
	IngestionService *service.IngestionService
	AnomalyService   *service.AnomalyService
}

func NewRouter(cfg config.Config, metrics *appmetrics.Collector, ingestionService *service.IngestionService, anomalyService *service.AnomalyService) http.Handler {
	deps := RouterDependencies{
		Config:           cfg,
		Metrics:          metrics,
		IngestionService: ingestionService,
		AnomalyService:   anomalyService,
	}

	tokenAuth := jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)
	handler := NewHandler(deps, tokenAuth)

	r := chi.NewRouter()
	r.Use(custommiddleware.CORS(cfg.AllowedOrigins))
	r.Use(custommiddleware.Prometheus(metrics))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/health", handler.Health)
	r.Handle("/metrics", metrics.Handler())
	r.Post("/auth/token", handler.IssueDemoToken)

	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))

			r.Post("/ingest", handler.CreateIngestionJob)
			r.Get("/anomalies", handler.ListAnomalies)
			r.Get("/billing-records", handler.ListBillingRecords)
			r.Get("/jobs", handler.ListJobs)
			r.Get("/metrics/summary", handler.MetricsSummary)
		})
	})

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		response := map[string]string{
			"name":        "Cost Guard API",
			"environment": cfg.Environment,
			"docs":        "See root README for architecture and API flow.",
		}
		writeJSON(w, http.StatusOK, response)
	})

	return r
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeText(w http.ResponseWriter, statusCode int, body string) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	w.WriteHeader(statusCode)
	_, _ = fmt.Fprint(w, body)
}
