rm -rf pkg/protobuf/*.go
protoc --go_out=./ --go-grpc_out=./ ./pkg/sdb-protobuf/*.proto --go-grpc_opt=require_unimplemented_servers=false
protoc --grpc-gateway_out . --grpc-gateway_opt logtostderr=true --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt grpc_api_configuration=./pkg/sdb-protobuf/sdb.yaml ./pkg/sdb-protobuf/*.proto
protoc --openapiv2_out . --openapiv2_opt grpc_api_configuration=./pkg/sdb-protobuf/sdb.yaml ./pkg/sdb-protobuf/sdb.proto