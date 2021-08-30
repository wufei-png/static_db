GO=go

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

VERSION := 1.0.0
BUILD := `git rev-parse --short HEAD`
TARGETS := elector
ALL_TARGETS:=$(TARGETS) $(TEST_TARGETS)
project=gitlab.sz.sensetime.com/viper/gocql-leaderelect

all: check build

build: $(TARGETS) $(TEST_TARGETS)

export CGO_LDFLAGS = $(LDFLAGS) -L$(PWD)/libs -L/opt/OpenBLAS/lib

$(TARGETS): $(SRC)
	$(GO) build -race $(project)/cmd/$@

$(TEST_TARGETS): $(SRC)
	$(GO) build $(project)/test/$@

.PHONY: clean all build check image

packages = $(shell go list ./...|grep -v /vendor/)
test: check
	$(GO) test -timeout=20m -race -v ${packages}

lint:
	@gometalinter --deadline=200s --config=.gometalint ./...

cov: check
	gocov test -timeout=20m -race -v $(packages) | gocov-html > coverage.html
	@cat coverage.html |grep "<code>Report Total</code>" | perl -nle 'print  "Total Coverage: $$1" if /([0-9.]+%)/g'

check:
	@$(GO) tool vet ${SRC}

clean:
	rm -f $(TARGETS) $(TEST_TARGETS)
