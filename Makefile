.DEFAULT_GOAL := help

SLACK_BOT_TOKEN := $(cat config/slack_bot_token)
SLACK_SIGNING_SECRET := $(cat config/slack_signing_secret)


build/%: src/%.go
#	@mkdir -p build
	go build -o $@ $<

.PHONY : compile
compile: ## compile
	go build -o build/main src/main.go

.PHONY : run
run: ## run project
	build/main

.PHONY : all
all: compile config/slack_bot_token config/slack_signing_secret
	export SLACK_BOT_TOKEN=$(shell cat ./config/slack_bot_token) && \
	export SLACK_SIGNING_SECRET=$(shell cat ./config/slack_signing_secret) && \
	export SLACK_REQUEST_CHANNEL=$(shell cat ./config/slack_requests_channel) && \
	build/main


.PHONY : gomod_tidy
gomod_tidy: ## run go mod tidy
	go mod tidy

.PHONY : gofmt
gofmt: ## run go fmt
	go fmt -x ./...

.PHONY : help
help: ## show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

