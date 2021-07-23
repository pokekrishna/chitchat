.EXPORT_ALL_VARIABLES:
COVERFILE := /tmp/cover.data
COVERHTML := /tmp/cover.html

.DEFAULT_GOAL := all

all: test

.PHONY: test
test:
	go test -coverprofile=$(COVERFILE) ./...
	@go tool cover -html=$(COVERFILE) -o $(COVERHTML)
	@echo ""
	@echo "Test coverage HTML published at $(COVERHTML)"

.PHONY: todo
todo:
	@echo "--- TODOs ---"
	@egrep -nir ".*//.*TODO.*:" . | grep -v "@egrep"
