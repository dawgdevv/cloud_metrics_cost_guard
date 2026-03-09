package main

import (
	"log"
	"net/http"

	"github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/config"
	apihttp "github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/http"
	"github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/ingestion"
	appmetrics "github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/metrics"
	"github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/service"
	"github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/store"
)

func main() {
	cfg := config.Load()
	repo := store.NewMemoryStore()
	metrics := appmetrics.NewCollector()
	detectionService := service.NewDetectionService()
	sources := map[string]ingestion.Source{
		"synthetic": ingestion.NewSyntheticSource(),
		"aws":       ingestion.NewAWSCostExplorerSource(cfg.AWSRegion, cfg.AWSAccountID),
	}
	ingestionService := service.NewIngestionService(repo, detectionService, metrics, sources)
	anomalyService := service.NewAnomalyService(repo)

	router := apihttp.NewRouter(cfg, metrics, ingestionService, anomalyService)

	log.Printf("api listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal(err)
	}
}
