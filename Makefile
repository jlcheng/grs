# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
GOGEN=$(GOCMD) generate
OUTDIR=out
MAIN_PRIME=cmd/grs/main.go
OUT_PRIME=out/grs

GENERATED=grs_stat_strings.go ui/gui_event_strings.go

VERSION=`head -n 1 VERSION.txt`

.PHONY: all
all: test build

.PHONY: build
build: $(OUT_PRIME)

.PHONY: test
test: $(GENERATED)
	$(GOTEST) -v jcheng/grs jcheng/grs/shexec

.PHONY: clean
clean: 
	rm -rf $(OUTDIR)

$(GENERATED): grs_stat.go
	$(GOGEN) .../ .../ui

install: all
	mv $(OUT_PRIME) $(HOME)/bin

run:
	$(GORUN) $(MAIN_PRIME)

.PHONY: fmt
fmt:
	gofmt -s -w ./
	golangci-lint run

out/grs: $(GENERATED)
	$(GOBUILD) -ldflags "-X main.Version=$(VERSION)" -o $(OUT_PRIME) $(MAIN_PRIME)
