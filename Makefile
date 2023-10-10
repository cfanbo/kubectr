PREFIX?=$(shell pwd)

NAME := kubectr
PKG := github.com/cfanbo/kubectr


BINDIR := ${PREFIX}/bin
GO111MODULE=on
CGO_ENABLED := 0
# Set any default go build tags
BUILDTAGS :=

# Add to compile time flags
VERSION := $(shell git describe --tags)
GITCOMMIT := $(shell git rev-parse --short HEAD)
BUILDMETA:=
GITUNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(GITUNTRACKEDCHANGES),)
	BUILDMETA := dirty
endif
CTIMEVAR=-X $(PKG)/internal/version.GitCommit=$(GITCOMMIT) \
	-X $(PKG)/internal/version.BuildMeta=$(BUILDMETA) \
	-X $(PKG)/internal/version.Version=$(VERSION) \
	-X $(PKG)/internal/version.ProjectURL=$(PKG)

GO ?= "go"
GO_LDFLAGS=-ldflags "-s -w $(CTIMEVAR)"
GOOSARCHES = linux/amd64 darwin/amd64 windows/amd64
GOOS = $(shell $(GO) env GOOS)
GOARCH= $(shell $(GO) env GOARCH)

.PHONY: all
all: clean lint test build ## Runs a clean, build, fmt, lint, test, and vet.

.PHONY: build
build:
	@echo "==> $@"
	@CGO_ENABLED=0 GO111MODULE=on $(GO) build -tags "$(BUILDTAGS)" ${GO_LDFLAGS} -o $(BINDIR)/$(NAME) ./cmd/main.go

.PHONY: buildtar
buildtar: clean lint test
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on $(GO) build -tags "$(BUILDTAGS)" ${GO_LDFLAGS} -o $(BINDIR)/$(NAME) ./cmd/main.go && tar -czf bin/kubectr_darwin_amd64.tar.gz LICENSE -C bin/ $(NAME)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 GO111MODULE=on $(GO) build -tags "$(BUILDTAGS)" ${GO_LDFLAGS} -o $(BINDIR)/$(NAME) ./cmd/main.go && tar -czf bin/kubectr_darwin_arm64.tar.gz LICENSE -C bin/ $(NAME)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on $(GO) build -tags "$(BUILDTAGS)" ${GO_LDFLAGS} -o $(BINDIR)/$(NAME) ./cmd/main.go && tar -czf bin/kubectr_linux_amd64.tar.gz LICENSE -C bin/ $(NAME)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 GO111MODULE=on $(GO) build -tags "$(BUILDTAGS)" ${GO_LDFLAGS} -o $(BINDIR)/$(NAME) ./cmd/main.go && tar -czf bin/kubectr_linux_arm64.tar.gz LICENSE -C bin/ $(NAME)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on $(GO) build -tags "$(BUILDTAGS)" ${GO_LDFLAGS} -o $(BINDIR)/$(NAME).exe ./cmd/main.go && zip -q bin/kubectr_windows_amd64.zip LICENSE -j bin/$(NAME).exe

.PHONY: test
test:
	@$(GO) test -v ./...

.PHONY: lint
lint: ## Verifies `golint` passes.
	@echo "==> $@"
	@golangci-lint run ./...

.PHONY: clean
clean:
	@echo "==> $@"
	$(RM) -r $(BINDIR)
