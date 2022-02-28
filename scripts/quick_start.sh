ps aux | grep "sdb" | grep 'config' | awk '{print "kill -9 " $2}' | sh -x
go run cmd/sdb/main.go -config ./configs/config.yml