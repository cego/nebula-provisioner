GRPC_VERSION = v1.1.0

models.pb.go: models.proto
	GOBIN="$(CURDIR)" go install google.golang.org/protobuf/cmd/protoc-gen-go
	PATH="$(CURDIR):$(PATH)" protoc --go_out=. --go_opt=paths=source_relative $<
	rm protoc-gen-go

agent-service.pb.go: agent-service.proto models.proto
	GOBIN="$(CURDIR)" go install google.golang.org/protobuf/cmd/protoc-gen-go
	GOBIN="$(CURDIR)" go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(GRPC_VERSION)
	PATH="$(CURDIR):$(PATH)" protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative $<
	rm protoc-gen-go protoc-gen-go-grpc

server-command.pb.go: server-command.proto models.proto
	GOBIN="$(CURDIR)" go install google.golang.org/protobuf/cmd/protoc-gen-go
	GOBIN="$(CURDIR)" go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(GRPC_VERSION)
	PATH="$(CURDIR):$(PATH)" protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative $<
	rm protoc-gen-go protoc-gen-go-grpc
