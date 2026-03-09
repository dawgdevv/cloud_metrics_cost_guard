API_DIR=apps/api
WEB_DIR=apps/web

.PHONY: api-run worker-run web-dev compose-up compose-down

api-run:
	cd $(API_DIR) && go run ./cmd/api

worker-run:
	cd $(API_DIR) && go run ./cmd/worker

web-dev:
	cd $(WEB_DIR) && npm install && npm run dev

compose-up:
	docker compose up --build

compose-down:
	docker compose down -v
