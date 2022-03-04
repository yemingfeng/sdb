package cli

import (
	"fmt"
	"github.com/abiosoft/ishell/v2"
	"github.com/yemingfeng/sdb/internal/pb"
	"github.com/yemingfeng/sdb/internal/util"
	"golang.org/x/net/context"
)

func RegisterGeoHashCmd(shell *ishell.Shell) {
	shell.AddCmd(newGHCreateCmd())
	shell.AddCmd(newGHDelCmd())
	shell.AddCmd(newGHAddCmd())
	shell.AddCmd(newGHPopCmd())
	shell.AddCmd(newGHGetBoxesCmd())
	shell.AddCmd(newGHGetNeighborsCmd())
	shell.AddCmd(newGHCountCmd())
	shell.AddCmd(newGHMembersCmd())
}

func newGHCreateCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "ghcreate",
		Help: "ghcreate key precision",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 2 {
				c.Println("args incorrect")
				return
			}
			precision, err := util.StringToInt32(c.Args[1])
			if err != nil {
				c.Println(err.Error())
				return
			}
			response, err := client.GHCreate(context.Background(), &pb.GHCreateRequest{Key: []byte(c.Args[0]), Precision: precision})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newGHDelCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "ghdel",
		Help: "ghdel key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.GHDel(context.Background(), &pb.GHDelRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newGHAddCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "ghadd",
		Help: "ghadd key id0 latitude0 longitude0 id1 latitude1 longitude1......",
		Func: func(c *ishell.Context) {
			if (len(c.Args)-1)%3 != 0 {
				c.Println("args incorrect")
				return
			}
			points := make([]*pb.Point, (len(c.Args)-1)/3)
			i := 0
			j := 1
			for i < (len(c.Args)-1)/3 {
				latitude, err := util.StringToDouble(c.Args[j+1])
				if err != nil {
					c.Println(err.Error())
					return
				}
				longitude, err := util.StringToDouble(c.Args[j+2])
				if err != nil {
					c.Println(err.Error())
					return
				}
				points[i] = &pb.Point{
					Id:        []byte(c.Args[j]),
					Latitude:  latitude,
					Longitude: longitude,
				}
				i += 1
				j += 3
			}
			response, err := client.GHAdd(context.Background(), &pb.GHAddRequest{Key: []byte(c.Args[0]), Points: points})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newGHPopCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "ghpop",
		Help: "ghpop key ids",
		Func: func(c *ishell.Context) {
			if len(c.Args) < 2 {
				c.Println("args incorrect")
				return
			}
			ids := make([][]byte, len(c.Args)-1)
			for i := 1; i < len(c.Args); i++ {
				ids[i-1] = []byte(c.Args[i])
			}
			response, err := client.GHPop(context.Background(), &pb.GHPopRequest{Key: []byte(c.Args[0]), Ids: ids})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Success)
			}
		},
	}
}

func newGHGetBoxesCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "ghgetboxes",
		Help: "ghgetboxes key latitude longitude",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 3 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			latitude, err := util.StringToDouble(c.Args[1])
			if err != nil {
				c.Println(err.Error())
				return
			}
			longitude, err := util.StringToDouble(c.Args[2])
			if err != nil {
				c.Println(err.Error())
				return
			}
			response, err := client.GHGetBoxes(context.Background(), &pb.GHGetBoxesRequest{Key: []byte(key), Latitude: latitude, Longitude: longitude})
			if err != nil {
				c.Println(err.Error())
			} else {
				for i := range response.Points {
					c.Println(string(response.Points[i].Id) + "\t" +
						fmt.Sprintf("%f", response.Points[i].Latitude) + "\t" +
						fmt.Sprintf("%f", response.Points[i].Longitude))
				}
			}
		},
	}
}

func newGHGetNeighborsCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "ghgetneighbors",
		Help: "ghgetneighbors key latitude longitude",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 3 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			latitude, err := util.StringToDouble(c.Args[1])
			if err != nil {
				c.Println(err.Error())
				return
			}
			longitude, err := util.StringToDouble(c.Args[2])
			if err != nil {
				c.Println(err.Error())
				return
			}
			response, err := client.GHGetNeighbors(context.Background(), &pb.GHGetNeighborsRequest{Key: []byte(key), Latitude: latitude, Longitude: longitude})
			if err != nil {
				c.Println(err.Error())
			} else {
				for i := range response.Points {
					c.Println(string(response.Points[i].Id) + "\t" +
						fmt.Sprintf("%f", response.Points[i].Latitude) + "\t" +
						fmt.Sprintf("%f", response.Points[i].Longitude) + "\t" +
						fmt.Sprintf("%d", response.Points[i].Distance))
				}
			}
		},
	}
}

func newGHCountCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "ghcount",
		Help: "ghcount key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.GHCount(context.Background(), &pb.GHCountRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				c.Println(response.Count)
			}
		},
	}
}

func newGHMembersCmd() *ishell.Cmd {
	return &ishell.Cmd{
		Name: "ghmembers",
		Help: "ghmembers key",
		Func: func(c *ishell.Context) {
			if len(c.Args) != 1 {
				c.Println("args incorrect")
				return
			}
			key := c.Args[0]
			response, err := client.GHMembers(context.Background(), &pb.GHMembersRequest{Key: []byte(key)})
			if err != nil {
				c.Println(err.Error())
			} else {
				for i := range response.Points {
					c.Println(string(response.Points[i].Id) + "\t" +
						fmt.Sprintf("%f", response.Points[i].Latitude) + "\t" +
						fmt.Sprintf("%f", response.Points[i].Longitude))
				}
			}
		},
	}
}
