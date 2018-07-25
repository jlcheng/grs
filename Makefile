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

PKGS=$(shell find src/jcheng/grs -maxdepth 1 -type d)

PRG2=cmd/grsnote/grsnote.go
PRG2_NAME=grsnote

all: test build
build: | $(OUTDIR) prg1 prg2

.PHONY: test
test: 
	$(GOTEST) -v jcheng/grs/test

vet:
	$(GOVET) $(PKGS)

clean: 
	$(GOCLEAN)
	rm -rf $(OUTDIR)

gofmt:
	$(GOFMT) -s -w $(PKGS)

install: all
	mv $(OUTDIR)/$(BINARY_NAME) $(HOME)/bin
	mv $(OUTDIR)/$(PRG2_NAME) $(HOME)/bin

run:
	$(GORUN) $(GOMAIN)

$(OUTDIR):
	mkdir -p $(OUTDIR)

prg1:
	$(GOBUILD) -o $(OUTDIR)/$(BINARY_NAME) $(GOMAIN)

prg2:
	$(GOBUILD) -o $(OUTDIR)/$(PRG2_NAME) $(PRG2)

