package conf

import (
	"flag"
	"github.com/yemingfeng/sdb/internal/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Store   Store   `yaml:"store"`
	Server  Server  `yaml:"server"`
	Cluster Cluster `yaml:"cluster"`
}

type Store struct {
	Engine string `yaml:"engine"`
	Path   string `yaml:"path"`
}

type Server struct {
	GRPCPort int `yaml:"grpc_port"`
	HttpPort int `yaml:"http_port"`
	Rate     int `yaml:"rate"`
}

type Cluster struct {
	NodeId  uint64 `yaml:"node_id"`
	Path    string `yaml:"path"`
	Address string `yaml:"address"`
	Master  string `yaml:"master"`
	Timeout int64  `yaml:"timeout"`
	Join    bool   `yaml:"join"`
}

var confLogger = util.GetLogger("conf")
var Conf Config

func init() {
	file := flag.String("config", "configs/config.yml", "config")

	confLogger.Printf("use config file: %s", *file)

	flag.Parse()

	bs, err := ioutil.ReadFile(*file)
	if err != nil {
		confLogger.Fatalf("read file %s %+v ", *file, err)
	}
	err = yaml.Unmarshal(bs, &Conf)
	if err != nil {
		confLogger.Fatalf("unmarshal: %+v", err)
	}

	confLogger.Printf("conf: %+v", Conf)
}
