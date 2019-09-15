export FYNE_THEME := light

GOPKG := github.com/dustinrohde/gifseekr
SRC := $(GOPATH)/src/$(GOPKG)
GOPKGS = $(GOPKG)/pkg/...
GOPKG_APP = $(GOPKG)/cmd/gifseekr

.PHONY: setup all build run test vet fmt lint check

all: build

build:
	go install $(if $V,-v) $(GOPKG_APP)

run:
	go run $(if $V,-v) $(GOPKG_APP)

test:
	go test $(if $V,-v) $(GOPKGS)

vet:
	go vet $(if $V,-v) $(GOPKGS) $(GOPKG_APP)

fmt:
	goimports -w $(if $V,-v) ./pkg/* ./cmd/*

lint:
	golint $(if $V,-v) $(GOPKGS) $(GOPKG_APP)

check: fmt lint test
