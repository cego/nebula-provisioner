BUILD_ARGS = -trimpath
LDFLAGS = -X main.Build=$(if $(BUILD_NUMBER),$(BUILD_NUMBER),"DEBUG")
CMD_PATH = "./cmd/"
# Set up OS specific bits
ifeq ($(OS),Windows_NT)
	GOVERSION := $(shell go version | awk "{print substr($$3, 3)}")
	GOISMIN := $(shell IF "$(GOVERSION)" GEQ "$(GOMINVERSION)" ECHO 1)
	CMD_SUFFIX = .exe
	NULL_FILE = nul
else
	GOVERSION := $(shell go version | awk '{print substr($$3, 3)}')
	GOISMIN := $(shell expr "$(GOVERSION)" ">=" "$(GOMINVERSION)")
	CMD_SUFFIX =
	NULL_FILE = /dev/null
endif

ALL_LINUX = linux-amd64 \
	linux-386

ALL = $(ALL_LINUX)

all: $(ALL:%=build/%/server) $(ALL:%=build/%/server-client) $(ALL:%=build/%/agent)

release: $(ALL:%=build/nebula-provisioner-%.tar.gz)

release-linux: $(ALL_LINUX:%=build/nebula-provisioner-%.tar.gz)

build/%/server: generate webapp .FORCE
	CGO_ENABLED=0 \
			GOOS=$(firstword $(subst -, , $*)) \
    		GOARCH=$(word 2, $(subst -, ,$*)) $(GOENV) \
			go build $(BUILD_ARGS) -o $@ -ldflags "$(LDFLAGS)" ${CMD_PATH}/server

build/%/server-client: .FORCE
	CGO_ENABLED=0 \
			GOOS=$(firstword $(subst -, , $*)) \
    		GOARCH=$(word 2, $(subst -, ,$*)) $(GOENV) \
			go build $(BUILD_ARGS) -o $@ -ldflags "$(LDFLAGS)" ${CMD_PATH}/server-client

build/%/agent: .FORCE
	CGO_ENABLED=0 \
			GOOS=$(firstword $(subst -, , $*)) \
    		GOARCH=$(word 2, $(subst -, ,$*)) $(GOENV) \
			go build $(BUILD_ARGS) -o $@ -ldflags "$(LDFLAGS)" ${CMD_PATH}/agent

build/nebula-provisioner-%.tar.gz: build/%/server build/%/server-client build/%/agent
	tar -zcv -C build/$* -f $@ server server-client agent

bin: generate webapp
	go build $(BUILD_ARGS) -ldflags "$(LDFLAGS)" $(BUILD_TAGS) -o ./bin/server${CMD_SUFFIX} ${CMD_PATH}/server
	go build $(BUILD_ARGS) -ldflags "$(LDFLAGS)" $(BUILD_TAGS) -o ./bin/server-client${CMD_SUFFIX} ${CMD_PATH}/server-client
	go build $(BUILD_ARGS) -ldflags "$(LDFLAGS)" $(BUILD_TAGS) -o ./bin/agent${CMD_SUFFIX} ${CMD_PATH}/agent

dev: BUILD_TAGS=-tags dev
dev: LDFLAGS=-X github.com/cego/nebula-provisioner/webapp.Dir=$(CURDIR)/webapp/
dev: bin

fmt:
	go fmt ./...

vet:
	go vet -v ./...

test:
	go test -v ./...

impsort:
	find . -iname '*.go' | grep -v '\.pb\.go$$' | xargs go run golang.org/x/tools/cmd/goimports -w

WEBAPP_DIRS = $(shell find ./webapp/ -type d  | grep -vE 'webapp/(node_modules|dist)')
WEBAPP_FILES = $(shell find ./webapp/ -type f -name '*' | grep -vE 'webapp/(node_modules|dist)')

webapp: ./webapp/ $(WEBAPP_DIRS) $(WEBAPP_FILES)
	$(MAKE) -C webapp build

.PHONY: generate
generate: protoc-gen-go protoc-gen-go-grpc buf gqlgen
	$(BUF) generate
	$(GQLGEN) generate

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
BUF = $(LOCALBIN)/buf
PROTOC_GEN_GO = $(LOCALBIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC = $(LOCALBIN)/protoc-gen-go-grpc
GQLGEN = $(LOCALBIN)/gqlgen

## Tool Versions
BUF_VERSION ?= v1.61.0
PROTOC_GEN_GO_VERSION ?= v1.36.10
PROTOC_GEN_GO_GRPC_VERSION ?= v1.6.0
GQLGEN_VERSION ?= v0.17.84

.PHONY: buf
buf: $(BUF) ## Download buf locally if necessary.
$(BUF): $(LOCALBIN)
	$(call go-install-tool,$(BUF),github.com/bufbuild/buf/cmd/buf,$(BUF_VERSION))

.PHONY: protoc-gen-go
protoc-gen-go: $(PROTOC_GEN_GO) ## Download protoc-gen-go locally if necessary.
$(PROTOC_GEN_GO): $(LOCALBIN)
	$(call go-install-tool,$(PROTOC_GEN_GO),google.golang.org/protobuf/cmd/protoc-gen-go,$(PROTOC_GEN_GO_VERSION))

.PHONY: protoc-gen-go-grpc
protoc-gen-go-grpc: $(PROTOC_GEN_GO_GRPC) ## Download protoc-gen-go-grpc locally if necessary.
$(PROTOC_GEN_GO_GRPC): $(LOCALBIN)
	$(call go-install-tool,$(PROTOC_GEN_GO_GRPC),google.golang.org/grpc/cmd/protoc-gen-go-grpc,$(PROTOC_GEN_GO_GRPC_VERSION))

.PHONY: gqlgen
gqlgen: $(GQLGEN) ## Download gqlgen locally if necessary.
$(GQLGEN): $(LOCALBIN)
	$(call go-install-tool,$(GQLGEN),github.com/99designs/gqlgen,$(GQLGEN_VERSION))


# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef

.FORCE:
.PHONY: test bin generate dev
.DEFAULT_GOAL := bin
