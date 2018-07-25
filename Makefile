# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOFMT=gofmt
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
GOMAIN=cmd/grs/grs.go
OUTDIR=out
# Sets GOPATH to the current project directory
export GOPATH=$(shell pwd)

WHAT := grs grsnote

all: test build
build: | $(OUTDIR)
	for target in $(WHAT); do \
		$(GOBUILD) -o $(OUTDIR)/$$target ./cmd/$$target; \
	done

.PHONY: test
test: 
	$(GOTEST) -v jcheng/grs/test

clean: 
	$(GOCLEAN)
	rm -rf $(OUTDIR)

gofmt:
	$(GOFMT) -s -w $(shell find src/jcheng/grs -maxdepth 1 -type d -not -name vendor)

install: all
	for target in $(WHAT); do \
		mv $(OUTDIR)/$$target $(HOME)/bin; \
	done

run:
	$(GORUN) $(GOMAIN)

$(OUTDIR):
	mkdir -p $(OUTDIR)
