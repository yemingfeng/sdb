rm -rf internal/pb/*.go
protoc --go_out=. --go-grpc_out=. ./internal/pb/protobuf-spec/*.proto --go-grpc_opt=require_unimplemented_servers=false
protoc --grpc-gateway_out . --grpc-gateway_opt logtostderr=true --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt grpc_api_configuration=./api/sdb.yaml ./internal/pb/protobuf-spec/*.proto
protoc --openapiv2_out . --openapiv2_opt grpc_api_configuration=./api/sdb.yaml ./internal/pb/protobuf-spec/sdb.proto