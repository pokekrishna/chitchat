.EXPORT_ALL_VARIABLES:
COVERFILE := /tmp/cover.data
COVERHTML := /tmp/cover.html

.DEFAULT_GOAL := all

all: test

test:
	go test -coverprofile=$(COVERFILE) -v ./...
	@go tool cover -html=$(COVERFILE) -o $(COVERHTML)
	@echo ""
	@echo "Test coverage HTML published at $(COVERHTML)"