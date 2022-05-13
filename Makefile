SRC = $(shell find . -name *.rs)
BINARY = pfi

.PHONY: help
help:
	@echo "pfi"
	@echo
	@echo "Usage:"
	@echo "  help:		show this help"
	@echo "  build:	build pfi"
	@echo "  dev-test: build and run against test image"

target/release/$(BINARY): $(SRC)
	cargo build --release

target/debug/$(BINARY): $(SRC)
	cargo build

.INTERMEDIATE: $(BINARY)
bin/$(BINARY): $(BINARY)
	mkdir -p bin
	mv $< $@

.PHONY: build-release
build-release: target/release/$(BINARY)
	cp target/release/$(BINARY) $(BINARY)
	@$(MAKE) --no-print-directory bin/$(BINARY)

.PHONY: build-debug
build-debug: target/debug/$(BINARY)
	cp target/debug/$(BINARY) $(BINARY)
	@$(MAKE) --no-print-directory bin/$(BINARY)
