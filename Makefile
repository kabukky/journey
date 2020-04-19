all: journey fmt vet lint

.PHONY: all fmt vet lint

PACKAGE = github.com/kabukky/journey
PKG_DIRS ?= authentication configuration conversion database date filenames flags helpers https plugins server slug structure templates watcher
PKG_FILES := $(foreach dir,$(PKG_DIRS), $(wildcard $(dir)/*.go))
MAIN_FILES = main.go

VET_LOG = vet.log
vet: vendor
	@$(foreach dir,$(PKG_DIRS),go vet $(VET_RULES) $(PACKAGE)/$(dir)|tee -a $(VET_LOG); )
	@go vet $(VET_RULES) $(MAIN_FILES) | tee -a $(VET_LOG)
	@[ ! -s $(VET_LOG) ]
	@rm $(VET_LOG)

GOLINT ?= golint
LINT_LOG = lint.log
lint: vendor
	@rm -f $(LINT_LOG)
	@$(foreach dir,$(PKG_DIRS),golint $(PACKAGE)/$(dir)|tee -a $(LINT_LOG); )
	@$(GOLINT) $(MAIN_FILES) | tee -a $(LINT_LOG)
	@[ ! -s $(LINT_LOG) ]
	@rm $(LINT_LOG)

GOIMPORTS ?= goimports
GOFMT ?= gofmt
fmt:
	@$(GOFMT) -s -w $(PKG_FILES)
	@$(GOIMPORTS) -w $(PKG_FILES)

journey: $(PKG_FILES) vendor
	go build
