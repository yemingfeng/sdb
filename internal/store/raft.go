package store

import (
	"errors"
	"fmt"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"github.com/yemingfeng/sdb/internal/conf"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var node *raft.Raft

func StartRaft() {
	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(conf.Conf.Cluster.NodeId)

	exist := true
	baseDir := filepath.Join(conf.Conf.Cluster.Path)
	if _, err := os.Stat(baseDir); err != nil {
		if os.IsNotExist(err) {
			exist = false
			if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
				log.Fatalln(err)
			}
		}
	}

	logStore, err := boltdb.NewBoltStore(filepath.Join(baseDir, "logs.dat"))
	if err != nil {
		log.Fatalln(err)
	}

	stableStore, err := boltdb.NewBoltStore(filepath.Join(baseDir, "stable.dat"))
	if err != nil {
		log.Fatalln(err)
	}

	snapshotStore, err := raft.NewFileSnapshotStore(baseDir, 3, os.Stderr)
	if err != nil {
		log.Fatalln(err)
	}

	addr, err := net.ResolveTCPAddr("tcp", conf.Conf.Cluster.Address)
	if err != nil {
		log.Fatalln(err)
	}

	transport, err := raft.NewTCPTransport(conf.Conf.Cluster.Address, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		log.Fatalln(err)
	}

	node, err = raft.NewRaft(raftConfig, NewFSM(), logStore, stableStore, snapshotStore, transport)
	if err != nil {
		log.Fatalln(err)
	}

	if !exist {
		if len(conf.Conf.Cluster.Master) == 0 {
			configuration := raft.Configuration{
				Servers: []raft.Server{
					{
						ID:      raftConfig.LocalID,
						Address: transport.LocalAddr(),
					},
				},
			}
			future := node.BootstrapCluster(configuration)
			if err := future.Error(); err != nil {
				log.Fatalln(err)
			}
		} else {
			resp, err := http.Get(fmt.Sprintf("http://%s/join?address=%s&nodeId=%s", conf.Conf.Cluster.Master, conf.Conf.Cluster.Address, conf.Conf.Cluster.NodeId))
			if err != nil {
				log.Fatalln(err)
			}
			defer func() {
				_ = resp.Body.Close()
			}()
		}
	}
}

func Apply(log *pb.Log) error {
	if node.State() != raft.Leader {
		return errors.New("only apply at master node")
	}
	if log == nil || len(log.LogEntries) == 0 {
		return nil
	}
	data, err := proto.Marshal(log)
	if err != nil {
		return err
	}
	return node.Apply(data, time.Duration(conf.Conf.Cluster.Timeout)*time.Millisecond).Error()
}

func HandleJoin(nodeId, address string) error {
	log.Printf("received join request for remote node %s at %s", nodeId, address)

	configFuture := node.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		log.Printf("failed to get raft configuration: %v", err)
		return err
	}

	for _, srv := range configFuture.Configuration().Servers {
		if srv.ID == raft.ServerID(nodeId) || srv.Address == raft.ServerAddress(address) {
			if srv.Address == raft.ServerAddress(address) && srv.ID == raft.ServerID(nodeId) {
				log.Printf("node %s at %s already member of cluster, ignoring join request", nodeId, address)
				return nil
			}

			future := node.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error removing existing node %s at %s: %s", nodeId, address, err)
			}
		}
	}

	f := node.AddVoter(raft.ServerID(nodeId), raft.ServerAddress(address), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}
	log.Printf("node %s at %s joined successfully", nodeId, address)
	return nil
}
