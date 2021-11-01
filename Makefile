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

build/%/server: .FORCE
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

bin: protocol server/store/store.pb.go impsort
	go build $(BUILD_ARGS) -ldflags "$(LDFLAGS)" -o ./bin/server${CMD_SUFFIX} ${CMD_PATH}/server
	go build $(BUILD_ARGS) -ldflags "$(LDFLAGS)" -o ./bin/server-client${CMD_SUFFIX} ${CMD_PATH}/server-client
	go build $(BUILD_ARGS) -ldflags "$(LDFLAGS)" -o ./bin/agent${CMD_SUFFIX} ${CMD_PATH}/agent

protocol: protocol/models.pb.go protocol/agent-service.pb.go protocol/server-command.pb.go

protocol/models.pb.go: protocol/models.proto
	$(MAKE) -C protocol models.pb.go

protocol/agent-service.pb.go: protocol/agent-service.proto
	$(MAKE) -C protocol agent-service.pb.go

protocol/server-command.pb.go: protocol/server-command.proto
	$(MAKE) -C protocol server-command.pb.go

server/store/store.pb.go: server/store/store.proto
	go build google.golang.org/protobuf/cmd/protoc-gen-go
	PATH="$(CURDIR):$(PATH)" protoc --go_out=. --go_opt=paths=source_relative $<
	rm protoc-gen-go

fmt:
	go fmt ./...

vet:
	go vet -v ./...

test:
	go test -v ./...

impsort:
	go build golang.org/x/tools/cmd/goimports
	find . -iname '*.go' | grep -v '\.pb\.go$$' | xargs ./goimports -w
	rm -f ./goimports

webapp:
	$(MAKE) -C webapp build

.FORCE:
.PHONY: test bin protocol webapp
.DEFAULT_GOAL := bin
