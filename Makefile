BUILD_ARGS = -trimpath
LDFLAGS = -X main.Build=$(BUILD_NUMBER)
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

build/%/server: webapp .FORCE
	GOOS=$(firstword $(subst -, , $*)) \
    		GOARCH=$(word 2, $(subst -, ,$*)) $(GOENV) \
			go build $(BUILD_ARGS) -o $@ -ldflags "$(LDFLAGS)" ${CMD_PATH}/server

build/%/server-client: .FORCE
	GOOS=$(firstword $(subst -, , $*)) \
    		GOARCH=$(word 2, $(subst -, ,$*)) $(GOENV) \
			go build $(BUILD_ARGS) -o $@ -ldflags "$(LDFLAGS)" ${CMD_PATH}/server-client

build/%/agent: .FORCE
	GOOS=$(firstword $(subst -, , $*)) \
    		GOARCH=$(word 2, $(subst -, ,$*)) $(GOENV) \
			go build $(BUILD_ARGS) -o $@ -ldflags "$(LDFLAGS)" ${CMD_PATH}/agent

build/nebula-provisioner-%.tar.gz: build/%/server build/%/server-client build/%/agent
	tar -zcv -C build/$* -f $@ server server-client agent

bin: protocol server/store/store.pb.go server/graph/generated/generated.go webapp
	go build $(BUILD_ARGS) -ldflags "$(LDFLAGS)" $(BUILD_TAGS) -o ./bin/server${CMD_SUFFIX} ${CMD_PATH}/server
	go build $(BUILD_ARGS) -ldflags "$(LDFLAGS)" $(BUILD_TAGS) -o ./bin/server-client${CMD_SUFFIX} ${CMD_PATH}/server-client
	go build $(BUILD_ARGS) -ldflags "$(LDFLAGS)" $(BUILD_TAGS) -o ./bin/agent${CMD_SUFFIX} ${CMD_PATH}/agent

dev: BUILD_TAGS=-tags dev
dev: LDFLAGS=-X github.com/slyngdk/nebula-provisioner/webapp.Dir=$(CURDIR)/webapp/
dev: bin

protocol: protocol/models.pb.go protocol/agent-service.pb.go protocol/server-command.pb.go
	go generate protocol/protocol.go

server/store/store.pb.go: server/store/store.proto
	go generate server/store/store.go

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

server/graph/generated/generated.go: server/graph/schema.graphqls gqlgen.yml
	go generate server/graph/resolver.go

.FORCE:
.PHONY: test bin protocol dev
.DEFAULT_GOAL := bin
