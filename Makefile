.PHONY: run-app
run-app:
	go run ./cmd/api/main.go --config_path=./config/ --config_file=config-local
