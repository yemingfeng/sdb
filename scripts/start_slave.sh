sh ./scripts/build_protobuf.sh

ps aux | grep "sdb" | grep 'slave' | awk '{print "kill -9 " $2}' | sh -x
go run cmd/sdb/sdb.go -config ./configs/slave1.yml
go run cmd/sdb/sdb.go -config ./configs/slave2.yml