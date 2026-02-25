GO=go
GOLANGCI_LINT=golangci-lint
GO_TEST=$(GO) test
GO_BUILD=$(GO) build
GO_CLEAN=$(GO) clean
GO_MOD=$(GO) mod

GO_BUILD_DIR=./build
SH_TOOLS_DIR=./tools/sh

SOURCES=$(shell find . -type f -name "*.go" -print)
SOURCE_DIRS=$(shell go list -f "{{.Dir}}" ./...)
BUILD_SHA:=$(shell git rev-parse --short HEAD 2>/dev/null || echo "0000000")
BUILD_VERSION:=$(shell cat ./VERSION 2>/dev/null || echo "0.0.0")
DIST_DIR=./dist

DOCKER_IMAGE ?= wile-extension-example
DOCKER_PLATFORM ?=
DOCKER_SHELL ?=

LDFLAGS=-ldflags "-X main.BuildSHA=$(BUILD_SHA) -X main.BuildVersion=$(BUILD_VERSION)"

# Detect host OS and architecture using Go conventions
HOST_OS := $(shell $(GO) env GOOS)
RAW_ARCH := $(shell uname -m)
ifeq ($(RAW_ARCH),x86_64)
HOST_ARCH := amd64
else
HOST_ARCH := $(RAW_ARCH)
endif

# Build all example binaries to ./dist/{os}/{arch}/.
# Rebuilds only when Go source files change.
#   make build
.PHONY: build
build:
	@mkdir -p $(DIST_DIR)/$(HOST_OS)/$(HOST_ARCH)
	@for dir in cmd/*/; do \
		name=$$(basename "$$dir"); \
		echo "Building $$name..."; \
		GOOS=$(HOST_OS) GOARCH=$(HOST_ARCH) $(GO_BUILD) -o $(DIST_DIR)/$(HOST_OS)/$(HOST_ARCH)/$$name $(LDFLAGS) ./$$dir; \
	done

# Build for a specific OS/architecture:
#   make build-darwin-arm64     # macOS Apple Silicon
#   make build-darwin-amd64     # macOS Intel
#   make build-linux-arm64      # Linux ARM64
#   make build-linux-amd64      # Linux x86-64
#   make build-all              # All OS/arch combinations
define BUILD_PLATFORM
.PHONY: build-$(1)-$(2)
build-$(1)-$(2):
	@mkdir -p $(DIST_DIR)/$(1)/$(2)
	@for dir in cmd/*/; do \
		name=$$$$(basename "$$$$dir"); \
		echo "Building $$$$name ($(1)/$(2))..."; \
		GOOS=$(1) GOARCH=$(2) $(GO_BUILD) -o $(DIST_DIR)/$(1)/$(2)/$$$$name $(LDFLAGS) ./$$$$dir; \
	done
endef

$(eval $(call BUILD_PLATFORM,darwin,arm64))
$(eval $(call BUILD_PLATFORM,darwin,amd64))
$(eval $(call BUILD_PLATFORM,linux,arm64))
$(eval $(call BUILD_PLATFORM,linux,amd64))

.PHONY: build-all
build-all: build-darwin-arm64 build-darwin-amd64 build-linux-arm64 build-linux-amd64

# Compile tests for all packages without running them.
# Useful for verifying that tests compile after refactoring.
#   make buildtest
.PHONY: buildtest
buildtest:
	for dir in $(SOURCE_DIRS); do \
	    if [ -d "$$dir" ]; then \
	        $(GO_TEST) -c -o /dev/null $$dir/...; \
	    fi \
	done

# ── CI: everything that must pass before merge ──────────────────────
# Set SKIP_LINT=1 when lint is handled externally (e.g., golangci-lint-action).
#   make ci
#   make ci SKIP_LINT=1
.PHONY: ci
ci: $(if $(SKIP_LINT),,lint) build test verify-mod
	@echo "CI passed"

# ── CD: release-specific validation ─────────────────────────────────
# Run before tagged releases. CI already passed on merge.
#   make cd
.PHONY: cd
cd: build test run-examples smoke-test
	@echo "CD passed"

# Run all tests with verbose output.
#   make test
.PHONY: test
test:
	$(GO_TEST) ./...

# Run all benchmarks with memory allocation statistics.
#   make bench
.PHONY: bench
bench:
	$(GO_TEST) -bench=. -benchmem ./...

# Run tests with coverage and print per-function coverage summary.
# Writes coverage profile to ./build/coverage.out.
#   make cover
.PHONY: cover
cover:
	@mkdir -p $(GO_BUILD_DIR)
	$(GO_TEST) -coverprofile=$(GO_BUILD_DIR)/coverage.out ./...
	$(GO) tool cover -func=$(GO_BUILD_DIR)/coverage.out

# Run tests with coverage and open an HTML report in the browser.
# Writes coverage profile to ./build/coverage.out and HTML to ./build/coverage.html.
#   make coverhtml
.PHONY: coverhtml
coverhtml:
	@mkdir -p $(GO_BUILD_DIR)
	$(GO_TEST) -coverprofile=$(GO_BUILD_DIR)/coverage.out ./...
	$(GO) tool cover -html=$(GO_BUILD_DIR)/coverage.out -o $(GO_BUILD_DIR)/coverage.html
	@echo "Coverage report: $(GO_BUILD_DIR)/coverage.html"
	open $(GO_BUILD_DIR)/coverage.html 2>/dev/null || xdg-open $(GO_BUILD_DIR)/coverage.html 2>/dev/null || echo "Open $(GO_BUILD_DIR)/coverage.html in your browser"

# Run tests with coverage and enforce per-package threshold (80%).
#   make covercheck
.PHONY: covercheck
covercheck:
	@mkdir -p $(GO_BUILD_DIR)
	$(GO_TEST) -coverprofile=$(GO_BUILD_DIR)/coverage.out ./...
	@bash $(SH_TOOLS_DIR)/covercheck.sh 80 $(GO_BUILD_DIR)/coverage.out

# Run golangci-lint on all packages.
#   make lint
.PHONY: lint
lint:
	$(GOLANGCI_LINT) -v run ./...

# Run golangci-lint with --fix to auto-correct fixable issues.
#   make fix
.PHONY: fix
fix:
	$(GOLANGCI_LINT) -v run --fix ./...

# Format all Go source files via golangci-lint.
#   make format
.PHONY: format
format:
	$(GOLANGCI_LINT) -v fmt -v ./...

# Tidy go.mod: add missing and remove unused dependencies.
#   make tidy
.PHONY: tidy
tidy:
	$(GO_MOD) tidy -e -x

# Remove all generated artifacts: build, test, module caches and output directories.
#   make clean
.PHONY: clean
clean: buildclean testclean modclean
	for dir in "$(DIST_DIR)" "$(GO_BUILD_DIR)"; do \
	    if [ -e "$$dir" ]; then rm -rvf "$$dir"; fi \
	done; \
	for dir in $(SOURCE_DIRS); do \
	    if [ -e "$$dir" ]; then find "$$dir" -name "*.test" -type f -exec rm -v \{\} \; ; fi \
	done

# Clear the Go build cache.
#   make buildclean
.PHONY: buildclean
buildclean:
	$(GO_CLEAN) -cache

# Clear the Go test and fuzz caches.
#   make testclean
.PHONY: testclean
testclean:
	$(GO_CLEAN) -testcache -fuzzcache

# Clear the Go module download cache.
#   make modclean
.PHONY: modclean
modclean:
	$(GO_CLEAN) -modcache

# Create an annotated git tag from the version in ./VERSION.
#   make tag
.PHONY: tag
tag:
	git tag -a $(BUILD_VERSION) -m "Release $(BUILD_VERSION)"
	@echo "Created tag $(BUILD_VERSION)"

# Bump the major version in VERSION (resets minor and patch to 0, preserves pre-release suffix).
#   make bump-major
#   v0.0.1 → v1.0.0
.PHONY: bump-major
bump-major:
	$(SH_TOOLS_DIR)/bump-version.sh major

# Bump the minor version in VERSION (resets patch to 0, preserves pre-release suffix).
#   make bump-minor
#   v0.0.1 → v0.1.0
.PHONY: bump-minor
bump-minor:
	$(SH_TOOLS_DIR)/bump-version.sh minor

# Bump the patch version in VERSION (preserves pre-release suffix).
#   make bump-patch
#   v0.0.1 → v0.0.2
.PHONY: bump-patch
bump-patch:
	$(SH_TOOLS_DIR)/bump-version.sh patch

# Verify go.sum integrity.
#   make verify-mod
.PHONY: verify-mod
verify-mod:
	$(GO_MOD) verify

# Run all example binaries and verify they exit 0.
#   make run-examples
.PHONY: run-examples
run-examples: build
	@$(SH_TOOLS_DIR)/run-examples.sh $(DIST_DIR)/$(HOST_OS)/$(HOST_ARCH)

# Smoke test: verify built binaries start and produce output.
#   make smoke-test
.PHONY: smoke-test
smoke-test: build
	@$(SH_TOOLS_DIR)/smoke-test.sh $(DIST_DIR)/$(HOST_OS)/$(HOST_ARCH)

# Build the Docker image.
#   make docker-build
#
# Cross-platform:
#   make docker-build DOCKER_PLATFORM=linux/amd64
.PHONY: docker-build
docker-build:
	DOCKER_IMAGE=$(DOCKER_IMAGE) DOCKER_PLATFORM=$(DOCKER_PLATFORM) $(SH_TOOLS_DIR)/docker-build.sh

# Open an interactive shell inside the Docker container.
#   make docker-shell
#   make docker-shell DOCKER_SHELL=/bin/sh
.PHONY: docker-shell
docker-shell:
	DOCKER_IMAGE=$(DOCKER_IMAGE) $(SH_TOOLS_DIR)/docker-shell.sh $(DOCKER_SHELL)
