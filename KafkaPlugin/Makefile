SHELL := /bin/bash
VERSION = $(shell cat ../VERSION)

all: clean docker_build

docker_build:
	@echo "---------------------------------------------------------------------------------"
	@echo "Starting docker build process, for KafkaPlugin......"
	@echo "---------------------------------------------------------------------------------"
	@echo ""
	@echo "Docker Build..."
	@echo "..............."
	@cp ../VERSION .
	@docker build -t kafkaplugin:1 . --network host
	@docker create --name kafkaplugin kafkaplugin:1
	@docker cp kafkaplugin:/KafkaPlugin/kafkaplugin-unittest-result.txt .
	@docker cp kafkaplugin:/KafkaPlugin/kafkaplugin-lint-result.txt .
	@docker cp kafkaplugin:/KafkaPlugin/coverage.html .
	@docker rm kafkaplugin
	@echo "docker build completed."

build:
	@echo Building KafkaPlugin
	@go mod download && go build -ldflags "-X main.appVersion=$(VERSION)" -o bin/kafkaplugin cmd/main.go || (echo "KafkaPlugin build failed"; exit 1)
	@echo Running go lint
	@go vet ./... > elasticsearchplugin-lint-result.txt
	@echo Build Successful.

build_package:
	@echo Creating KafkaPlugin package
	@mkdir -p bin
	@docker create --name kafkaplugin kafkaplugin:1
	@docker cp kafkaplugin:/KafkaPlugin/bin/kafkaplugin ./bin/
	@docker rm kafkaplugin
	@mkdir -p package && mkdir -p package/bin && cp -R resources package && cp bin/kafkaplugin package/bin && chmod 777 package/bin/kafkaplugin
	@cd package && zip -r KafkaPlugin-$(VERSION).zip bin resources && rm -rf bin resources
	@echo Package created at package/KafkaPlugin-$(VERSION).zip

test:
	@echo "Started :KafkaPlugin Tests"
	@echo "---> Starting ElasticSearch"
	@runuser -l elasticsearch -c 'cd /kafkaplugin/elasticsearch-7.0.1/bin && nohup bash -c "./elasticsearch -E http.port=9299 -E transport.tcp.port=9399 &" && sleep 20'
	@echo Running Tests
	@go test ./... -v -coverprofile=coverage.out > kafkaplugin-unittest-result.txt || echo "KafkaPlugin unit tests failed."
	@echo Running go coverage
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Completed :KafkaPlugin Tests"

clean:
	@echo "Started :KafkaPlugin CleanUp"
	@rm -rf package bin/kafkaplugin VERSION kafkaplugin-unittest-result.txt kafkaplugin-lint-result.txt
	@echo "Completed :KafkaPlugin CleanUp"
