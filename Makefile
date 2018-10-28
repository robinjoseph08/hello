DIRS ?= $(shell find . -name '*.go' | grep --invert-match 'vendor' | xargs -n 1 dirname | sort --unique)
SETUP_PKGS := \
	github.com/alecthomas/gometalinter \
	github.com/golang/dep/cmd/dep \

BFLAGS ?=
LFLAGS ?=
TFLAGS ?=

COVERAGE_PROFILE ?= coverage.out

default: install

.PHONY: clean
clean:
	@echo "---> Cleaning"
	rm -rf ./vendor

.PHONY: html
html:
	@echo "---> Generating HTML coverage report"
	go tool cover -html $(COVERAGE_PROFILE)

.PHONY: install
install:
	@echo "---> Installing dependencies"
	dep ensure

.PHONY: lint
lint:
	@echo "---> Linting"
	gometalinter --vendor --tests $(LFLAGS) $(DIRS)

.PHONY: setup
setup:
	@echo "--> Installing linter"
	go get -u -v $(SETUP_PKGS)
	gometalinter --install

.PHONY: test
test:
	@echo "---> Testing"
	ENVIRONMENT=test go test ./... -race -coverprofile $(COVERAGE_PROFILE) $(TFLAGS)
