all: journey

.PHONY: all fmt clean

PACKAGE = github.com/rkuris/journey
PKG_DIRS ?= authentication configuration conversion database date filenames flags \
	    helpers https notifications plugins server slug structure templates watcher
PKGS := $(foreach dir,$(PKG_DIRS), $(PACKAGE)/$(dir))
PKG_FILES := $(foreach dir,$(PKG_DIRS), $(wildcard $(dir)/*.go))
MAIN_FILES = main.go

VET_LOG = vet.log
vet: vendor
	@$(foreach dir,$(PKG_DIRS),go vet $(VET_RULES) $(PACKAGE)/$(dir)|tee -a $(VET_LOG); )
	@go vet $(VET_RULES) $(MAIN_FILES) | tee -a $(VET_LOG)
	@[ ! -s $(VET_LOG) ]
	@rm $(VET_LOG)

GOIMPORTS ?= go imports
GOFMT ?= go fmt
fmt:
	@$(GOFMT) $(PKGS)

journey: $(PKG_FILES) vendor
	go build

clean:
	rm -f journey lint.log vet.log


vendor:
	mkdir vendor
