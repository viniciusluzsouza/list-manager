.PHONY: help

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
SHELL:=/bin/bash

all: help

.PHONY:
help:
	@echo "List Manager API"
	@echo ""
	@echo "run                       Starts application with dependent containers"
	@echo "clean                     Remove all application related containers"
	@echo "generate-mocks            Generate mocks for unit testing"
	@echo ""

.PHONY:
build: clean
	@docker-compose build --force-rm 

.PHONY:
run: build
	@docker-compose up -d

.PHONY:
clean:
	@docker-compose down

.PHONY:
generate-mocks:
	@rm -f ./internal/modules/**/mocks/*
	@go generate ./internal/** ./internal/modules/**
