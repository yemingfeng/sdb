package util

import (
	"github.com/bwmarrin/snowflake"
	"log"
)

var node *snowflake.Node

func init() {
	node2, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatal("generate snowflake node error", err)
	}
	node = node2
}

func GetOrderingKey() int64 {
	return node.Generate().Int64()
}
