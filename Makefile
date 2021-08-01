BUILD_ARGS = -trimpath
LDFLAGS = -X main.Build=$(BUILD_NUMBER)
CMD_PATH = "./cmd/"

# Set up OS specific bits
ifeq ($(OS),Windows_NT)
	#TODO: we should be able to ditch awk as well
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

bin: protocol
	go build $(BUILD_ARGS) -ldflags "$(LDFLAGS)" -o ./bin/server${CMD_SUFFIX} ${CMD_PATH}/server
	go build $(BUILD_ARGS) -ldflags "$(LDFLAGS)" -o ./bin/server-client${CMD_SUFFIX} ${CMD_PATH}/server-client
	go build $(BUILD_ARGS) -ldflags "$(LDFLAGS)" -o ./bin/agent${CMD_SUFFIX} ${CMD_PATH}/agent

protocol: protocol/agent-service.pb.go protocol/server-command.pb.go

protocol/agent-service.pb.go: protocol/agent-service.proto
	$(MAKE) -C protocol agent-service.pb.go

protocol/server-command.pb.go: protocol/server-command.proto
	$(MAKE) -C protocol server-command.pb.go

vet:
	go vet -v ./...

test:
	go test -v ./...

.FORCE:
.PHONY: test bin protocol
.DEFAULT_GOAL := bin
