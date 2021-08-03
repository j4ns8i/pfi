SRC = $(shell find . -name '*.go')
BINARY = pfi

.PHONY: help
help:
	@echo "pfi"
	@echo
	@echo "Usage:"
	@echo "  help:		show this help"
	@echo "  build:	build pfi"

$(BINARY): $(SRC)
	go build -o $(BINARY)

.PHONY: build
build: $(BINARY)
