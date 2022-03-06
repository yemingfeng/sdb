package cli

import (
	"github.com/abiosoft/ishell/v2"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"github.com/yemingfeng/sdb/internal/util"
	"golang.org/x/net/context"
)

func RegisterBloomFilterCmd(shell *ishell.Shell) {
	shell.AddCmd(newBFCreateCmd())
	shell.AddCmd(newBFDelCmd())
	shell.AddCmd(newBFAddCmd())
	shell.AddCmd(newBFExistCmd())
}

func newBFCreateCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "bfcreate",
		Help: "bfcreate key n p",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 3 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			n, err := util.StringToUInt32(c.Args[1])
			if err != nil {
				c.Println(err.Error())
				return
			}
			p, err := util.StringToDouble(c.Args[2])
			if err != nil {
				c.Println(err.Error())
				return
			}
			response, err := client.BFCreate(context.Background(), &pb.BFCreateRequest{Key: []byte(key), N: n, P: p})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newBFDelCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "bfdel",
		Help: "bfdel key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.BFDel(context.Background(), &pb.BFDelRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newBFAddCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "bfadd",
		Help: "bfadd key values",
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
			response, err := client.BFAdd(context.Background(), &pb.BFAddRequest{Key: []byte(key), Values: bsValues})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newBFExistCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "bfexist",
		Help: "bfexist key values",
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
			response, err := client.BFExist(context.Background(), &pb.BFExistRequest{Key: []byte(key), Values: bsValues})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Exists)
			}
		},
	}
}
