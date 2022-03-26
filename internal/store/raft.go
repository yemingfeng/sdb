package store

import (
	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
	"github.com/yemingfeng/sdb/internal/conf"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
	"log"
	"path/filepath"
	"time"
)

var clusterId = uint64(1)
var node *dragonboat.NodeHost

func StartRaft() {
	initialMembers := make(map[uint64]string)
	if len(conf.Conf.Cluster.Master) == 0 {
		initialMembers[conf.Conf.Cluster.NodeId] = conf.Conf.Cluster.Address
	} else {
		initialMembers[1] = conf.Conf.Cluster.Master
	}

	rc := config.Config{
		NodeID:               conf.Conf.Cluster.NodeId,
		ClusterID:            clusterId,
		ElectionRTT:          10,
		HeartbeatRTT:         1,
		CheckQuorum:          true,
		SnapshotEntries:      10,
		CompactionOverhead:   5,
		EntryCompressionType: config.Snappy,
	}

	path := filepath.Join(conf.Conf.Cluster.Path)
	nhc := config.NodeHostConfig{
		WALDir:         path,
		NodeHostDir:    path,
		RTTMillisecond: 200,
		RaftAddress:    conf.Conf.Cluster.Address}
	var err error
	node, err = dragonboat.NewNodeHost(nhc)
	if err != nil {
		log.Fatalln(err)
	}
	if err := node.StartOnDiskCluster(initialMembers, len(conf.Conf.Cluster.Master) != 0, NewFSM, rc); err != nil {
		log.Fatalln(err)
	}
}

func Apply(pbLog *pb.Log) error {
	if pbLog == nil || len(pbLog.LogEntries) == 0 {
		return nil
	}
	data, err := proto.Marshal(pbLog)
	if err != nil {
		log.Printf("error on marshal pbLog: [%+v], err: [%+v]", pbLog, err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.Conf.Cluster.Timeout)*time.Millisecond)
	defer cancel()
	result, err := node.SyncPropose(ctx, node.GetNoOPSession(clusterId), data)
	if err != nil {
		log.Printf("error on: [%s], result: [%+v], err: [%+v]", pbLog, result, err)
	}
	return err
}
