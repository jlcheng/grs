# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
GOMAIN=cmd/grs/grs.go
BINARY_NAME=grs
OUTDIR=out

ifeq ($(OS),Windows_NT)
	BINARY_NAME=grs.exe
endif


all: test build
build: | $(OUTDIR)
	$(GOBUILD) -o $(OUTDIR)/$(BINARY_NAME) $(GOMAIN)

.PHONY: test
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -rf $(OUTDIR)

install: all
	mv $(OUTDIR)/$(BINARY_NAME) $(HOME)/bin

run:
	$(GORUN) $(GOMAIN)

$(OUTDIR):
	mkdir -p $(OUTDIR)
