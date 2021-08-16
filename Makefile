SRC = $(shell find . -name '*.go')
BINARY = pfi

.PHONY: help
help:
	@echo "pfi"
	@echo
	@echo "Usage:"
	@echo "  help:		show this help"
	@echo "  build:	build pfi"
	@echo "  dev-test: build and run against test image"

$(BINARY): $(SRC)
	go build -o $(BINARY)

.PHONY: build
build: $(BINARY)

.PHONY: dev-test
dev-test: test-output.yaml.tmpl test.jpg $(BINARY)
	./$(BINARY) test-output.yaml.tmpl test-output.yaml test.jpg
