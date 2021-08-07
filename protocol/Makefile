GO111MODULE = on
export GO111MODULE

models.pb.go: models.proto
	go build google.golang.org/protobuf/cmd/protoc-gen-go
	PATH="$(CURDIR):$(PATH)" protoc --go_out=. --go_opt=paths=source_relative $<
	rm protoc-gen-go

agent-service.pb.go: agent-service.proto models.proto
	go mod download google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go build google.golang.org/protobuf/cmd/protoc-gen-go
	go build google.golang.org/grpc/cmd/protoc-gen-go-grpc
	PATH="$(CURDIR):$(PATH)" protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative $<
	rm protoc-gen-go protoc-gen-go-grpc

server-command.pb.go: server-command.proto models.proto
	go mod download google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go build google.golang.org/protobuf/cmd/protoc-gen-go
	go build google.golang.org/grpc/cmd/protoc-gen-go-grpc
	PATH="$(CURDIR):$(PATH)" protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative $<
	rm protoc-gen-go protoc-gen-go-grpc
