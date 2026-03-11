.PHONY: proto
proto: ## Generate Go code by buf
	@echo "Linting and generating code with buf..."
	buf lint
	buf generate
	@echo "Done!"

.PHONY: clean
clean: ## Clean up generated files
	find . -name "*.pb.go" -type f -delete
	find . -name "*.pb.go" -type f -delete

.PHONY: master
master: ## Run Master server
	go run cmd/master/main.go

worker: ## Run Worker server
	go run cmd/worker/main.go

.PHONY: help
help: ## Show list of available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

