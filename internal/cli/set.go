package cli

import (
	"github.com/abiosoft/ishell/v2"
	"github.com/yemingfeng/sdb/internal/pb"
	"golang.org/x/net/context"
)

func RegisterSetCmd(shell *ishell.Shell) {
	shell.AddCmd(newSPushCmd())
	shell.AddCmd(newSPopCmd())
	shell.AddCmd(newSExistCmd())
	shell.AddCmd(newSDelCmd())
	shell.AddCmd(newSCountCmd())
	shell.AddCmd(newSMembersCmd())
}

func newSPushCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "spush",
		Help: "spush key values",
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
			response, err := client.SPush(context.Background(), &pb.SPushRequest{Key: []byte(key), Values: bsValues})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newSPopCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "spop",
		Help: "spop key values",
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
			response, err := client.SPop(context.Background(), &pb.SPopRequest{Key: []byte(key), Values: bsValues})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newSExistCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "sexist",
		Help: "sexist key values",
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
			response, err := client.SExist(context.Background(), &pb.SExistRequest{Key: []byte(key), Values: bsValues})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Exists)
			}
		},
	}
}

func newSDelCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "sdel",
		Help: "sdel key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.SDel(context.Background(), &pb.SDelRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newSCountCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "scount",
		Help: "scount key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.SCount(context.Background(), &pb.SCountRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Count)
			}
		},
	}
}

func newSMembersCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "smembers",
		Help: "smembers key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.SMembers(context.Background(), &pb.SMembersRequest{Key: []byte(key)})
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
