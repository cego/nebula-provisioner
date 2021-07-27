GO111MODULE = on
export GO111MODULE

nebula-provisioner.pb.go: nebula-provisioner.proto
	go build google.golang.org/protobuf/cmd/protoc-gen-go
	go build google.golang.org/grpc/cmd/protoc-gen-go-grpc
	PATH="$(CURDIR):$(PATH)" protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative $<
	rm protoc-gen-go protoc-gen-go-grpc

server-command.pb.go: server-command.proto
	go build google.golang.org/protobuf/cmd/protoc-gen-go
	go build google.golang.org/grpc/cmd/protoc-gen-go-grpc
	PATH="$(CURDIR):$(PATH)" protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative $<
	rm protoc-gen-go protoc-gen-go-grpc
