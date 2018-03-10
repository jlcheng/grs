# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
GOMAIN=cmd/grs/grs.go
BINARY_NAME=grs

ifeq ($(OS),Windows_NT)
	BINARY_NAME=grs.exe
endif


all: test build
build: 
	$(GOBUILD) -o $(BINARY_NAME) $(GOMAIN)

.PHONY: test
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

install: all
	mv $(BINARY_NAME) $(HOME)/bin

run:
	$(GORUN) $(GOMAIN)
