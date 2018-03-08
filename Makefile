# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
GOMAIN=cmd/v2/grs2.go
BINARY_NAME=grs2

ifeq ($(OS),Windows_NT)
	BINARY_NAME=grs2.exe
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

install: build
	mv $(BINARY_NAME) $(HOME)/bin

run:
	$(GORUN) $(GOMAIN)
