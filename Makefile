# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOFMT=gofmt
GOVET=$(GOCMD) vet
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
GOMAIN=cmd/grs/grs.go
BINARY_NAME=grs
OUTDIR=out
# Sets GOPATH to the current project directory
export GOPATH=$(shell pwd)


WHAT := grs grsdbh

all: test build

.PHONY: test
test: 
	$(GOTEST) -v jcheng/grs/test

vet:
	$(GOVET) jcheng/grs/compat jcheng/grs/config jcheng/grs/display jcheng/grs/gittest jcheng/grs/grs jcheng/grs/grsdb jcheng/grs/grsio jcheng/grs/script jcheng/grs/status jcheng/grs/test

clean: 
	$(GOCLEAN)
	rm -rf $(OUTDIR)

fmt:
	$(GOFMT) -s -w `find src/jcheng/grs -maxdepth 1 -and -type d -and -not -name vendor`

install: all
	for target in $(WHAT); do \
		mv $(OUTDIR)/$$target $(HOME)/bin; \
	done

run:
	$(GORUN) $(GOMAIN)

$(OUTDIR):
	mkdir -p $(OUTDIR)

build:
	for target in $(WHAT); do \
		$(GOBUILD) -o $(OUTDIR)/$$target ./cmd/$$target; \
	done

