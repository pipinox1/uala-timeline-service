# Makefile
app_name := uala-timeline-service
version ?= latest

.PHONY: build

build:
	docker build -t $(app_name):$(version) .