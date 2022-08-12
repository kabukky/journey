all: journey fmt vet

.PHONY: all fmt vet clean

PACKAGE = github.com/rkuris/journey
PKG_DIRS ?= authentication configuration conversion database date filenames flags \
	    helpers https notifications plugins server slug structure templates watcher
PKG_FILES := $(foreach dir,$(PKG_DIRS), $(wildcard $(dir)/*.go))
MAIN_FILES = main.go

VET_LOG = vet.log
vet: vendor
	@$(foreach dir,$(PKG_DIRS),go vet $(VET_RULES) $(PACKAGE)/$(dir)|tee -a $(VET_LOG); )
	@go vet $(VET_RULES) $(MAIN_FILES) | tee -a $(VET_LOG)
	@[ ! -s $(VET_LOG) ]
	@rm $(VET_LOG)

GOIMPORTS ?= goimports
GOFMT ?= gofmt
fmt:
	@$(GOFMT) -s -w $(PKG_FILES)
	@$(GOIMPORTS) -w $(PKG_FILES)

journey: $(PKG_FILES) vendor
	go build

clean:
	rm -f journey lint.log vet.log


vendor:
	mkdir vendor
