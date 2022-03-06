package cli

import (
	"github.com/abiosoft/ishell/v2"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"github.com/yemingfeng/sdb/internal/util"
	"golang.org/x/net/context"
)

func RegisterBitsetCmd(shell *ishell.Shell) {
	shell.AddCmd(newBSCreateCmd())
	shell.AddCmd(newBSDelCmd())
	shell.AddCmd(newBSSetRangeCmd())
	shell.AddCmd(newBSMSetCmd())
	shell.AddCmd(newBSGetRangeCmd())
	shell.AddCmd(newBSMGetCmd())
	shell.AddCmd(newBSCountCmd())
	shell.AddCmd(newBSCountRangeCmd())
}

func newBSCreateCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "bscreate",
		Help: "bscreate key size",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 2 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			size, err := util.StringToUInt32(c.Args[1])
			if err != nil {
				c.Println(err.Error())
				return
			}
			response, err := client.BSCreate(context.Background(), &pb.BSCreateRequest{Key: []byte(key), Size: size})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newBSDelCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "bsdel",
		Help: "bsdel key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.BSDel(context.Background(), &pb.BSDelRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newBSSetRangeCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "bssetrange",
		Help: "bssetrange key start end value",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 4 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			start, err := util.StringToUInt32(c.Args[1])
			if err != nil {
				c.Println(err.Error())
			}
			end, err := util.StringToUInt32(c.Args[2])
			if err != nil {
				c.Println(err.Error())
			}
			value, err := util.StringToBoolean(c.Args[3])
			if err != nil {
				c.Println(err.Error())
			}
			response, err := client.BSSetRange(context.Background(), &pb.BSSetRangeRequest{Key: []byte(key), Start: start, End: end, Value: value})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newBSMSetCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "bsmset",
		Help: "bsmset key bits value",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 3 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			argBits := c.Args[1:len(c.Args)]
			bits := make([]uint32, len(argBits))
			for i := range argBits {
				bit, err := util.StringToUInt32(argBits[i])
				if err != nil {
					c.Println(err.Error())
					return
				}
				bits[i] = bit
			}

			value, err := util.StringToBoolean(c.Args[len(c.Args)-1])
			if err != nil {
				c.Println(err.Error())
			}
			response, err := client.BSMSet(context.Background(), &pb.BSMSetRequest{Key: []byte(key), Bits: bits, Value: value})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newBSGetRangeCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "bsgetrange",
		Help: "bsgetrange key start end",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 3 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			start, err := util.StringToUInt32(c.Args[1])
			if err != nil {
				c.Println(err.Error())
			}
			end, err := util.StringToUInt32(c.Args[2])
			if err != nil {
				c.Println(err.Error())
			}
			response, err := client.BSGetRange(context.Background(), &pb.BSGetRangeRequest{Key: []byte(key), Start: start, End: end})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Values)
			}
		},
	}
}

func newBSMGetCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "bsmget",
		Help: "bsmget key bits",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 2 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			argBits := c.Args[1:len(c.Args)]
			bits := make([]uint32, len(argBits))
			for i := range argBits {
				bit, err := util.StringToUInt32(argBits[i])
				if err != nil {
					c.Println(err.Error())
					return
				}
				bits[i] = bit
			}
			response, err := client.BSMGet(context.Background(), &pb.BSMGetRequest{Key: []byte(key), Bits: bits})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Values)
			}
		},
	}
}

func newBSCountCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "bsmcount",
		Help: "bsmcount key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]

			response, err := client.BSCount(context.Background(), &pb.BSCountRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Count)
			}
		},
	}
}

func newBSCountRangeCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "bscountrange",
		Help: "bscountrange key start end",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 3 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			start, err := util.StringToUInt32(c.Args[1])
			if err != nil {
				c.Println(err.Error())
			}
			end, err := util.StringToUInt32(c.Args[2])
			if err != nil {
				c.Println(err.Error())
			}
			response, err := client.BSCountRange(context.Background(), &pb.BSCountRangeRequest{Key: []byte(key), Start: start, End: end})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Count)
			}
		},
	}
}
