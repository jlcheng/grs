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


PRG2=cmd/grsnote/grsnote.go
PRG2_NAME=grsnote

all: test build
build: | $(OUTDIR) prg1 prg2

.PHONY: test
test: 
	$(GOTEST) -v jcheng/grs/test
	$(GOFMT) -l .

vet:
	$(GOVET) jcheng/grs/compat jcheng/grs/config jcheng/grs/display jcheng/grs/gittest jcheng/grs/grs jcheng/grs/grsdb jcheng/grs/grsio jcheng/grs/script jcheng/grs/status jcheng/grs/test

clean: 
	$(GOCLEAN)
	rm -rf $(OUTDIR)

gofmt:
	$(GOFMT) -s -w .

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

