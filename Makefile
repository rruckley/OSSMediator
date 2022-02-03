SHELL := /bin/bash

all:
	@cd OSSMediatorCollector && make all
	@cd ElasticsearchPlugin && make all
	@cd KafkaPlugin && make all
