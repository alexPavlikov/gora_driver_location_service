.PHONY: run-app
run-app:
	go run ./cmd/api/main.go --config_path=./config/ --config_file=config-local

.PHONY: run-kafka
run-kafka:
	.\bin\windows\zookeeper-server-start.bat .\config\zookeeper.properties
	.\bin\windows\kafka-server-start.bat .\config\server.properties
