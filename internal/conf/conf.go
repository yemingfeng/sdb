package conf

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Store  Store  `yaml:"store"`
	Server Server `yaml:"server"`
}

type Store struct {
	Engine string `yaml:"engine"`
	Path   string `yaml:"path"`
}

type Server struct {
	GRPCPort           int   `yaml:"grpc_port"`
	HttpPort           int   `yaml:"http_port"`
	Rate               int   `yaml:"rate"`
	SlowQueryThreshold int64 `yaml:"slow_query_threshold"`
}

var Conf Config

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)

	file := flag.String("config", "configs/config.yml", "config")

	log.Printf("use config file: %s", *file)

	flag.Parse()

	bs, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf("read file %s %+v ", *file, err)
	}
	err = yaml.Unmarshal(bs, &Conf)
	if err != nil {
		log.Fatalf("unmarshal: %+v", err)
	}

	log.Printf("conf: %+v", Conf)
}
