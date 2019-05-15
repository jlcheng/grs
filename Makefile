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

GENERATED=script/grs_stat_strings.go

.PHONY: all
all: test build

.PHONY: build
build: $(OUT_PRIME)

.PHONY: test
test: $(GENERATED)
	$(GOTEST) -v jcheng/grs/script jcheng/grs/shexec

.PHONY: clean
clean: 
	rm -rf $(OUTDIR)

$(GENERATED): script/grs_stat.go
	$(GOGEN) .../script

install: all
	mv $(OUT_PRIME) $(HOME)/bin

run:
	$(GORUN) $(MAIN_PRIME)

.PHONY: fmt
fmt:
	gofmt -w ./

out/grs: $(GENERATED)
	$(GOBUILD) -o $(OUT_PRIME) $(MAIN_PRIME)
