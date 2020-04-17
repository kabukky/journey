ALL_SRC := $(shell find . -name "*.go" | grep -v -e vendor \
	-e ".*/\..*" \
	-e ".*/_.*")

GOIMPORTS=goimports

all: journey fmt vet lint

.PHONY: all fmt vet lint

PKG_DIRS ?= authentication configuration conversion database date filenames flags helpers https plugins server slug structure templates watcher
PKG_FILES := $(foreach dir,$(PKG_DIRS), $(wildcard $(dir)/*.go))

VET_LOG = vet.log
vet: vendor
	@$(foreach dir,$(PKG_DIRS),go vet $(VET_RULES) github.com/kabukky/journey/$(dir)|tee -a $(VET_LOG); )
	@go vet $(VET_RULES) main.go | tee -a $(VET_LOG)
	@[ ! -s $(VET_LOG) ]
	@rm $(VET_LOG)

LINT_LOG = lint.log
lint: vendor
	@rm -f $(LINT_LOG)
#	@$(foreach dir,$(PKG_DIRS),golint github.com/kabukky/journey/$(dir)|tee -a $(LINT_LOG); )
	@golint main.go | tee -a $(LINT_LOG)
	@[ ! -s $(LINT_LOG) ]
	@rm $(LINT_LOG)

fmt:
	@gofmt -s -w $(ALL_SRC)
	@$(GOIMPORTS) -w $(ALL_SRC)

journey: $(PKG_FILES) vendor
	go build
