SHELL := /bin/bash
VERSION = $(shell cat ../VERSION)

all: clean docker_build

docker_build:
	@echo "---------------------------------------------------------------------------------"
	@echo "Starting docker build process, for OSSMediatorCollector......"
	@echo "---------------------------------------------------------------------------------"
	@echo ""
	@echo "Docker Build..."
	@echo "..............."
	@cp ../VERSION .
	@docker build -t ossmediatorcollector:1 . --network host
	@docker create --name ossmediatorcollector ossmediatorcollector:1
	@docker cp ossmediatorcollector:/OSSMediatorCollector/collector-unittest-result.txt .
	@docker cp ossmediatorcollector:/OSSMediatorCollector/collector-lint-result.txt .
	@docker cp ossmediatorcollector:/OSSMediatorCollector/coverage.html .
	@docker rm ossmediatorcollector
	@echo "docker build completed."

build:
	@echo Building OSSMediatorCollector
	@go mod download && go build -ldflags "-X main.appVersion=$(VERSION)" -o bin/collector cmd/collector.go || (echo "OSSMediatorCollector build failed"; exit 1)
	@echo Running go lint
	@go vet ./... > collector-lint-result.txt
	@echo Build Successful.

build_storesecret:
	@echo Building storesecret
	@go mod download && go build -o bin/storesecret storesecret/storesecret.go || (echo "storesecret build failed"; exit 1)
	@echo storesecret build successful.

build_package:
	@echo Creating OSSMediatorCollector package
	@mkdir -p bin
	@docker create --name ossmediatorcollector ossmediatorcollector:1
	@docker cp ossmediatorcollector:/OSSMediatorCollector/bin/collector ./bin/
	@docker cp ossmediatorcollector:/OSSMediatorCollector/bin/storesecret ./bin/
	@docker rm ossmediatorcollector
	@mkdir -p package && mkdir -p package/bin && cp -R resources package && cp bin/* package/bin/ && chmod 777 package/bin/*
	@cd package && zip -r OSSMediatorCollector-$(VERSION).zip bin resources && rm -rf bin resources
	@echo Package created at package/OSSMediatorCollector-$(VERSION).zip

test:
	@echo "Started :OSSMediatorCollector Tests"
	@echo Running Tests
	@go test ./... -coverprofile=coverage.out -v > collector-unittest-result.txt || echo "OSSMediatorCollector unit tests failed."
	@echo Running go coverage
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Completed :OSSMediatorCollector Tests"

clean:
	@echo "Started :OSSMediatorCollector CleanUp"
	@rm -rf package bin/collector bin/storesecret coverage.html VERSION collector-unittest-result.txt collector-lint-result.txt
	@echo "Completed :OSSMediatorCollector CleanUp"
