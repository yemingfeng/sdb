package util

import (
	"github.com/bwmarrin/snowflake"
)

var keyLogger = GetLogger("key")
var node *snowflake.Node

func init() {
	node2, err := snowflake.NewNode(1)
	if err != nil {
		keyLogger.Fatal("generate snowflake node error", err)
	}
	node = node2
}

func GetOrderingKey() int64 {
	return node.Generate().Int64()
}
