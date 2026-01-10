.PHONY: help fetch-spec generate build install test clean

help:
	@echo "TrueNAS Terraform Provider - Development Commands"
	@echo ""
	@echo "  make fetch-spec   - Download latest API spec from TrueNAS"
	@echo "  make generate     - Generate provider code from spec"
	@echo "  make build        - Build provider binary"
	@echo "  make install      - Install provider locally"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean generated files"
	@echo ""
	@echo "Quick workflow:"
	@echo "  make fetch-spec generate build install"

fetch-spec:
	@echo "Fetching latest API specification..."
	@if [ -z "$$TRUENAS_HOST" ] || [ -z "$$TRUENAS_TOKEN" ]; then \
		echo "ERROR: Set TRUENAS_HOST and TRUENAS_TOKEN environment variables"; \
		echo "Example: export TRUENAS_HOST=192.168.1.100"; \
		echo "         export TRUENAS_TOKEN=your-api-token"; \
		exit 1; \
	fi
	python3 fetch_methods.py

generate:
	@echo "Generating provider code..."
	python3 generate.py

build:
	@echo "Building provider..."
	go build

install: build
	@echo "Installing provider locally..."
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/bmanojlovic/truenas/0.1.0/linux_amd64/
	cp terraform-provider-truenas ~/.terraform.d/plugins/registry.terraform.io/bmanojlovic/truenas/0.1.0/linux_amd64/
	@echo "✅ Provider installed"

test:
	go test -v ./internal/...

clean:
	rm -f terraform-provider-truenas
	rm -f truenas-methods-*.json
