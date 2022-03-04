package cli

import (
	"fmt"
	"github.com/abiosoft/ishell/v2"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/util"
	"golang.org/x/net/context"
)

func RegisterSortedSetCmd(shell *ishell.Shell) {
	shell.AddCmd(newZPushCmd())
	shell.AddCmd(newZPopCmd())
	shell.AddCmd(newZRangeCmd())
	shell.AddCmd(newZExistCmd())
	shell.AddCmd(newZDelCmd())
	shell.AddCmd(newZCountCmd())
	shell.AddCmd(newZMembersCmd())
}

func newZPushCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "zpush",
		Help: "zpush key value0 score0 key1 score1......",
		Func: func(c *ishell.Context) {
			if len(c.Args)%2 != 1 {
				c.Println("args incorrect")
				return
			}
			tuples := make([]*pb.Tuple, len(c.Args)/2)
			i := 0
			j := 1
			for i < len(c.Args)/2 {
				score, err := util.StringToDouble(c.Args[j+1])
				if err != nil {
					c.Println(err.Error())
					return
				}
				tuples[i] = &pb.Tuple{Value: []byte(c.Args[j]), Score: score}
				i += 1
				j += 2
			}
			response, err := client.ZPush(context.Background(), &pb.ZPushRequest{Key: []byte(c.Args[0]), Tuples: tuples})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newZPopCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "zpop",
		Help: "zpop key values",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 2 {
				c.Println("args incorrect")
				return
			}
			values := make([][]byte, len(c.Args)-1)
			for i := 1; i < len(c.Args); i++ {
				values[i-1] = []byte(c.Args[i])
			}
			response, err := client.ZPop(context.Background(), &pb.ZPopRequest{Key: []byte(c.Args[0]), Values: values})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newZRangeCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "zrange",
		Help: "zrange key offset limit",
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
			response, err := client.ZRange(context.Background(), &pb.ZRangeRequest{Key: []byte(key), Offset: offset, Limit: limit})
			if err != nil {
				c.Println(err.Error())
			} else {
				for i := range response.Tuples {
					c.Println(string(response.Tuples[i].Value) + "\t" + fmt.Sprintf("%32.32f", response.Tuples[i].Score))
				}
			}
		},
	}
}

func newZExistCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "zexist",
		Help: "zexist key values",
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
			response, err := client.ZExist(context.Background(), &pb.ZExistRequest{Key: []byte(key), Values: bsValues})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Exists)
			}
		},
	}
}

func newZDelCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "zdel",
		Help: "zdel key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.ZDel(context.Background(), &pb.ZDelRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newZCountCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "zcount",
		Help: "zcount key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.ZCount(context.Background(), &pb.ZCountRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Count)
			}
		},
	}
}

func newZMembersCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "zmembers",
		Help: "zmembers key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.ZMembers(context.Background(), &pb.ZMembersRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				for i := range response.Tuples {
					c.Println(string(response.Tuples[i].Value) + "\t" + fmt.Sprintf("%32.32f", response.Tuples[i].Score))
				}
			}
		},
	}
}
