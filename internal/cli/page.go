package cli

import (
	"github.com/abiosoft/ishell/v2"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"github.com/yemingfeng/sdb/internal/util"
	"golang.org/x/net/context"
)

func RegisterPageCmd(shell *ishell.Shell) {
	shell.AddCmd(newPListCmd())
}

func newPListCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "plist",
		Help: "plist dataType key offset limit",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 3 && len(c.Args) != 4 {
				c.Println("args incorrect")
				return
			}
			dataType, err := util.StringToUInt32(c.Args[0])
			if err != nil {
				c.Println(err.Error())
				return
			}
			var key string
			var offset int32
			var limit uint32
			if len(c.Args) == 3 {
				offset, err = util.StringToInt32(c.Args[1])
				if err != nil {
					c.Println(err.Error())
					return
				}
				limit, err = util.StringToUInt32(c.Args[2])
				if err != nil {
					c.Println(err.Error())
					return
				}
			} else {
				key = c.Args[1]
				offset, err = util.StringToInt32(c.Args[2])
				if err != nil {
					c.Println(err.Error())
					return
				}
				limit, err = util.StringToUInt32(c.Args[3])
				if err != nil {
					c.Println(err.Error())
					return
				}
			}
			response, err := client.PList(context.Background(), &pb.PListRequest{DataType: pb.DataType(dataType), Key: []byte(key), Offset: offset, Limit: limit})
			if err != nil {
				c.Println(err.Error())
			} else {
				strKeys := make([]string, len(response.Keys))
				for i := range response.Keys {
					strKeys[i] = string(response.Keys[i])
				}
				c.Println(strKeys)
			}
		},
	}
}
