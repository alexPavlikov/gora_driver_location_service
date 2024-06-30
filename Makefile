.PHONY: build-app
build-app:
	go run ./cmd/api/main.go --config=./config/config-local.yaml
