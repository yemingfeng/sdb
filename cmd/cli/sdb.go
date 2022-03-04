package main

import (
	"github.com/abiosoft/ishell/v2"
	"github.com/yemingfeng/sdb/internal/cli"
)

func main() {
	shell := ishell.New()
	shell.Println("sdb cli")
	cli.RegisterBitsetCmd(shell)
	cli.RegisterPageCmd(shell)
	cli.RegisterBloomFilterCmd(shell)
	cli.RegisterHyperLogLogCmd(shell)
	cli.RegisterListCmd(shell)
	cli.RegisterSetCmd(shell)
	cli.RegisterStringCmd(shell)
	cli.RegisterMapCmd(shell)
	cli.RegisterSortedSetCmd(shell)
	cli.RegisterGeoHashCmd(shell)
	shell.Run()
}
