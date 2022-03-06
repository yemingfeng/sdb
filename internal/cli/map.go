package cli

import (
	"github.com/abiosoft/ishell/v2"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
)

func RegisterMapCmd(shell *ishell.Shell) {
	shell.AddCmd(newMPushCmd())
	shell.AddCmd(newMPopCmd())
	shell.AddCmd(newMExistCmd())
	shell.AddCmd(newMDelCmd())
	shell.AddCmd(newMCountCmd())
	shell.AddCmd(newMMembersCmd())
}

func newMPushCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "mpush",
		Help: "mpush key key0 value0 key1 value1......",
		Func: func(c *ishell.Context) {
			if (len(c.Args)-1)%2 != 0 {
				c.Println("args incorrect")
				return
			}
			pairs := make([]*pb.Pair, (len(c.Args)-1)/2)
			i := 0
			j := 1
			for i < len(c.Args)/2 {
				pairs[i] = &pb.Pair{Key: []byte(c.Args[j]), Value: []byte(c.Args[j+1])}
				i += 1
				j += 2
			}
			response, err := client.MPush(context.Background(), &pb.MPushRequest{Key: []byte(c.Args[0]), Pairs: pairs})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newMPopCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "mpop",
		Help: "mpop key keys",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 2 {
				c.Println("args incorrect")
				return
			}
			keys := make([][]byte, len(c.Args)-1)
			for i := 1; i < len(c.Args); i++ {
				keys[i-1] = []byte(c.Args[i])
			}
			response, err := client.MPop(context.Background(), &pb.MPopRequest{Key: []byte(c.Args[0]), Keys: keys})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newMExistCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "mexist",
		Help: "mexist key keys",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 2 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			keys := c.Args[1:len(c.Args)]
			bsKeys := make([][]byte, len(keys))
			for i := range keys {
				bsKeys[i] = []byte(keys[i])
			}
			response, err := client.MExist(context.Background(), &pb.MExistRequest{Key: []byte(key), Keys: bsKeys})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Exists)
			}
		},
	}
}

func newMDelCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "mdel",
		Help: "mdel key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.MDel(context.Background(), &pb.MDelRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newMCountCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "mcount",
		Help: "mcount key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.MCount(context.Background(), &pb.MCountRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Count)
			}
		},
	}
}

func newMMembersCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "mmembers",
		Help: "mmembers key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.MMembers(context.Background(), &pb.MMembersRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				for i := range response.Pairs {
					c.Println(string(response.Pairs[i].Key) + "\t" + string(response.Pairs[i].Value))
				}
			}
		},
	}
}
