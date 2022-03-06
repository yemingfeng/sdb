package cli

import (
	"github.com/abiosoft/ishell/v2"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
)

func RegisterHyperLogLogCmd(shell *ishell.Shell) {
	shell.AddCmd(newHLLCreateCmd())
	shell.AddCmd(newHLLDelCmd())
	shell.AddCmd(newHLLAddCmd())
	shell.AddCmd(newHLLCountCmd())
}

func newHLLCreateCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "hllcreate",
		Help: "hllcreate key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.HLLCreate(context.Background(), &pb.HLLCreateRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newHLLDelCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "hlldel",
		Help: "hlldel key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.HLLDel(context.Background(), &pb.HLLDelRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newHLLAddCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "hlladd",
		Help: "hlladd key values",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 2 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			values := c.Args[1:len(c.Args)]
			bsValues := make([][]byte, len(values))
			for i := range values {
				bsValues[i] = []byte(values[i])
			}
			response, err := client.HLLAdd(context.Background(), &pb.HLLAddRequest{Key: []byte(key), Values: bsValues})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newHLLCountCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "hllcount",
		Help: "hllcount key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.HLLCount(context.Background(), &pb.HLLCountRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Count)
			}
		},
	}
}
