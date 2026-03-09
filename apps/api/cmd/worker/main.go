package main

import (
	"log"
	"time"

	"github.com/nishant-raj/multi_cloud_optimizer/apps/api/internal/config"
)

func main() {
	cfg := config.Load()
	log.Printf("worker scaffold started for environment=%s nats=%s", cfg.Environment, cfg.NATSURL)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Printf("worker heartbeat: plug JetStream consumer and anomaly processing here")
	}
}
