package util

import (
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node
var idLogger = GetLogger("id")

func init() {
	newNode, err := snowflake.NewNode(1)
	if err != nil {
		idLogger.Fatalf("can not new snowflake node, err: %+v", err)
	}
	node = newNode
}

func NextId() uint64 {
	return uint64(node.Generate())
}
