rm -rf pkg/protobuf/*.go
protoc --go_out=./ --go-grpc_out=./ ./internal/protobuf/*.proto