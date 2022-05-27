SRC = $(shell find . -name *.rs)
BINARY = pfi

.PHONY: help
help: ## Show this help
	@grep -E '^[a-z.A-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: install
install: ## Install pfi
	cargo install --path .

.PHONY: uninstall
uninstall: ## Uninstall pfi
	cargo uninstall

target/release/$(BINARY): $(SRC)
	cargo build --release

target/debug/$(BINARY): $(SRC)
	cargo build

.INTERMEDIATE: $(BINARY)
bin/$(BINARY): $(BINARY)
	mkdir -p bin
	mv $< $@

.PHONY: build-release
build-release: target/release/$(BINARY) ## Build release version
	cp target/release/$(BINARY) $(BINARY)
	@$(MAKE) --no-print-directory bin/$(BINARY)

.PHONY: build-debug
build-debug: target/debug/$(BINARY) ## Build debug version
	cp target/debug/$(BINARY) $(BINARY)
	@$(MAKE) --no-print-directory bin/$(BINARY)

.PHONY: test
test: ## Run unit tests
	cargo test
