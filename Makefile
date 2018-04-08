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

PRG2=cmd/grsnote/grsnote.go
PRG2_NAME=grsnote

ifeq ($(OS),Windows_NT)
	BINARY_NAME=grs.exe
	PRG2_NAME=grsnote.exe
endif


all: test build
build: | $(OUTDIR) prg1 prg2

.PHONY: test
test: 
	$(GOTEST) -v ./test
clean: 
	$(GOCLEAN)
	rm -rf $(OUTDIR)

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
