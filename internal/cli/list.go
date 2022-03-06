package cli

import (
	"github.com/abiosoft/ishell/v2"
	"github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
)

func RegisterListCmd(shell *ishell.Shell) {
	shell.AddCmd(newLRPushCmd())
	shell.AddCmd(newLLPushCmd())
	shell.AddCmd(newLPopCmd())
	shell.AddCmd(newLRangeCmd())
	shell.AddCmd(newLExistCmd())
	shell.AddCmd(newLDelCmd())
	shell.AddCmd(newLCountCmd())
	shell.AddCmd(newLMembersCmd())
}

func newLRPushCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "lrpush",
		Help: "lrpush key values",
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
			response, err := client.LRPush(context.Background(), &pb.LRPushRequest{Key: []byte(key), Values: bsValues})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newLLPushCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "llpush",
		Help: "llpush key values",
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
			response, err := client.LLPush(context.Background(), &pb.LLPushRequest{Key: []byte(key), Values: bsValues})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newLPopCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "lpop",
		Help: "lpop key values",
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
			response, err := client.LPop(context.Background(), &pb.LPopRequest{Key: []byte(key), Values: bsValues})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newLRangeCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "lrange",
		Help: "lrange key offset limit",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 3 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			offset, err := util.StringToInt32(c.Args[1])
			if err != nil {
				c.Println(err.Error())
			}
			limit, err := util.StringToUInt32(c.Args[2])
			if err != nil {
				c.Println(err.Error())
			}
			response, err := client.LRange(context.Background(), &pb.LRangeRequest{Key: []byte(key), Offset: offset, Limit: limit})
			if err != nil {
				c.Println(err.Error())
			} else {
				strKeys := make([]string, len(response.Values))
				for i := range response.Values {
					strKeys[i] = string(response.Values[i])
				}
				c.Println(strKeys)
			}
		},
	}
}

func newLExistCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "lexist",
		Help: "lexist key values",
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
			response, err := client.LExist(context.Background(), &pb.LExistRequest{Key: []byte(key), Values: bsValues})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Exists)
			}
		},
	}
}

func newLDelCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "ldel",
		Help: "ldel key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.LDel(context.Background(), &pb.LDelRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newLCountCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "lcount",
		Help: "lcount key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.LCount(context.Background(), &pb.LCountRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Count)
			}
		},
	}
}

func newLMembersCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "lmembers",
		Help: "lmembers key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.LMembers(context.Background(), &pb.LMembersRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				for i := range response.Values {
					c.Println(string(response.Values[i]))
				}
			}
		},
	}
}
