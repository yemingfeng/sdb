package cli

import (
	"github.com/abiosoft/ishell/v2"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/util"
	"golang.org/x/net/context"
)

func RegisterStringCmd(shell *ishell.Shell) {
	shell.AddCmd(newSetCmd())
	shell.AddCmd(newMSetCmd())
	shell.AddCmd(newSetNXCmd())
	shell.AddCmd(newGetCmd())
	shell.AddCmd(newMGetCmd())
	shell.AddCmd(newDelCmd())
	shell.AddCmd(newIncrCmd())
}

func newSetCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "set",
		Help: "set key value",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 2 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			value := c.Args[1]
			response, err := client.Set(context.Background(), &pb.SetRequest{Key: []byte(key), Value: []byte(value)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newMSetCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "mset",
		Help: "mset key0 value0 key1 value1 ......",
		Func: func(c *ishell.Context) {
			if len(c.Args)%2 != 0 {
				c.Println("args incorrect")
				return
			}
			keys := make([][]byte, len(c.Args)/2)
			values := make([][]byte, len(c.Args)/2)
			i := 0
			j := 0
			for i < len(c.Args)/2 {
				keys[i] = []byte(c.Args[j])
				values[i] = []byte(c.Args[j+1])
				i += 1
				j += 2
			}
			response, err := client.MSet(context.Background(), &pb.MSetRequest{Keys: keys, Values: values})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newSetNXCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "setnx",
		Help: "setnx key value",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 2 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			value := c.Args[1]
			response, err := client.SetNX(context.Background(), &pb.SetNXRequest{Key: []byte(key), Value: []byte(value)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newGetCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "get",
		Help: "get key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.Get(context.Background(), &pb.GetRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(string(response.Value))
			}
		},
	}
}

func newMGetCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "mget",
		Help: "mget keys",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 2 {
				c.Println("args incorrect")
				return
			}
			keys := make([][]byte, len(c.Args))
			for i := 0; i < len(c.Args); i++ {
				keys[i] = []byte(c.Args[i])
			}
			response, err := client.MGet(context.Background(), &pb.MGetRequest{Keys: keys})
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

func newDelCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "del",
		Help: "del key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.Del(context.Background(), &pb.DelRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newIncrCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "incr",
		Help: "incr key delta",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 2 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			delta, err := util.StringToInt32(c.Args[1])
			if err != nil {
				c.Println(err.Error())
				return
			}
			response, err := client.Incr(context.Background(), &pb.IncrRequest{Key: []byte(key), Delta: delta})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}
