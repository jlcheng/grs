# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
OUTDIR=out
MAIN_PRIME=cmd/grs/main.go
OUT_PRIME=out/grs

.PHONY: all
all: test build

.PHONY: build
build: $(OUT_PRIME)

.PHONY: test
test: 
	$(GOTEST) -v jcheng/grs/...

.PHONY: clean
clean: 
	$(GOCLEAN)
	rm -rf $(OUTDIR)

install: all
	mv $(OUT_PRIME) $(HOME)/bin

run:
	$(GORUN) $(MAIN_PRIME)

out/grs:
	$(GOBUILD) -o $(OUT_PRIME) $(MAIN_PRIME)
