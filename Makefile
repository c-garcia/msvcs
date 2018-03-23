.PHONY: clean e2e test test-cover wip dep \
	notif notif-linux notif-container \
	infra \
	help \
	users-double users-double-linux users-double-container \
	start-all stop-all ps-all logs-notif start-devel stop-devel ps-devel

.DEFAULT_GOAL := help


RMQ_URL ?= amqp://admin:admin@localhost:5672/
USER_SVC_URL ?= http://localhost:9999/api/
NOTIF_SVC_URL ?= http://localhost:10000/api/
MH_URL ?= http://localhost:8025/api/
SMTP_HOST ?= localhost
SMTP_PORT ?= 1025
USERS_BIN ?= users-double
NOTIF_BIN ?= notif

export RMQ_URL
export USER_SVC_URL
export NOTIF_SVC_URL
export MH_URL
export SMTP_HOST
export SMTP_PORT


wip: ## Run tests marked as WIP
	go test -tags=wip ./...

test: ## Run unit and component tests
	go test ./...

test-cover: ## Run unit and component tests with coverage enabled
	go test -cover ./...

e2e: ## Run e2e tests
	go test -tags=e2e ./...

notif: ## Build notif binary
	go build -o bin/$(NOTIF_BIN) cmd/notif/notif.go

notif-linux: ## Build Linux version of the notif binary
	NOTIF_BIN=notif.amd64 GOARCH=amd64 GOOS=linux make notif

notif-container: ## Build notif container
notif-container: notif-linux
	cp bin/notif.amd64 docker/obliquo/notif
	cd docker/obliquo/notif && docker build -t obliquo/notif .

users-double: ## Build users service double
	go build -o bin/$(USERS_BIN) cmd/users_double/users_double.go

users-double-linux: ## Build Linux version of the users service double
	USERS_BIN=users-double.amd64 GOARCH=amd64 GOOS=linux make users-double

users-double-container: ## Build users service container
users-double-container: users-double-linux
	cp bin/users-double.amd64 docker/obliquo/users-double
	cd docker/obliquo/users-double && docker build -t obliquo/users-double .

clean: ## Clean
	find . -type f -name "*~" -exec rm -f \{\} \;
	rm -f bin/*
	go clean

dep: ## Update go vendoring
	dep ensure

infra: ## Infrastructure containers
	cd docker/obliquo/clustered-rabbitmq && docker build -t obliquo/clustered-rabbitmq .
	cd docker/obliquo/rabbitmq-lb && docker build -t obliquo/rabbitmq-lb .

start-test: ## Start test environment
start-test: users-double-container notif-container infra
	docker-compose -f docker-compose.yml -f notif-compose.yml up

stop-test: ## Stop test environment
	docker-compose -f docker-compose.yml -f notif-compose.yml down

ps-test: ## Status of the test environment
	docker-compose -f docker-compose.yml -f notif-compose.yml ps

logs-notif-test: ## Notifier logs in the test environment
	docker-compose -f docker-compose.yml -f notif-compose.yml logs -f notif

start-devel: ## Start devel environment
start-devel: users-double-container infra
	docker-compose up

stop-devel: ## Stop devel environment
	docker-compose down

ps-devel: ## Status of the devel environment
	docker-compose -f docker-compose.yml -f notif-compose.yml ps

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


