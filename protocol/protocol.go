package protocol

//go:generate protoc --go_out=. --go_opt=paths=source_relative models.proto
//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative agent-service.proto models.proto
//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative server-command.proto models.proto
