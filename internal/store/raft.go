package store

import (
	"fmt"
	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
	"github.com/lni/dragonboat/v3/logger"
	"github.com/yemingfeng/sdb/internal/conf"
	"github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

var raftLogger = util.GetLogger("raft")
var clusterId = uint64(1)
var node *dragonboat.NodeHost

func StartCluster() {
	rc := config.Config{
		NodeID:               conf.Conf.Cluster.NodeId,
		ClusterID:            clusterId,
		ElectionRTT:          10,
		HeartbeatRTT:         1,
		CheckQuorum:          false,
		SnapshotEntries:      50,
		CompactionOverhead:   5,
		EntryCompressionType: config.Snappy}

	path := filepath.Join(conf.Conf.Cluster.Path)
	nhc := config.NodeHostConfig{
		WALDir:         path,
		NodeHostDir:    path,
		RTTMillisecond: 200,
		RaftAddress:    conf.Conf.Cluster.Address}
	var err error
	logger.GetLogger("rsm").SetLevel(logger.WARNING)
	node, err = dragonboat.NewNodeHost(nhc)
	if err != nil {
		raftLogger.Fatalln(err)
	}
	if len(conf.Conf.Cluster.Master) == 0 {
		initialMembers := map[uint64]string{
			conf.Conf.Cluster.NodeId: conf.Conf.Cluster.Address,
		}
		if err := node.StartOnDiskCluster(initialMembers, false, NewFSM, rc); err != nil {
			raftLogger.Fatalln(err)
		}
	} else {
		if err := node.StartOnDiskCluster(nil, true, NewFSM, rc); err != nil {
			raftLogger.Fatalln(err)
		}
		if conf.Conf.Cluster.Join {
			resp, err := http.Get(fmt.Sprintf("http://%s/join?nodeId=%d&address=%s", conf.Conf.Cluster.Master, conf.Conf.Cluster.NodeId, conf.Conf.Cluster.Address))
			if err != nil {
				raftLogger.Fatalln(err)
			}
			bs, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				raftLogger.Fatalln(err)
			}
			if "ok" != string(bs) {
				raftLogger.Fatalln("join failed")
			}
			defer func() {
				_ = resp.Body.Close()
			}()
		}
	}
}

func StopCluster() {
	if node != nil {
		node.Stop()
	}
}

func Apply(pbLog *pb.Log) error {
	if pbLog == nil || len(pbLog.LogEntries) == 0 {
		return nil
	}
	data, err := proto.Marshal(pbLog)
	if err != nil {
		raftLogger.Printf("error on marshal pbLog: [%+v], err: [%+v]", pbLog, err)
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.Conf.Cluster.Timeout)*time.Millisecond)
	defer cancel()
	result, err := node.SyncPropose(ctx, node.GetNoOPSession(clusterId), data)
	if err != nil {
		raftLogger.Printf("error on: [%s], result: [%+v], err: [%+v]", pbLog, result, err)
	}
	return err
}

func HandleJoin(nodeId uint64, address string) error {
	raftLogger.Printf("received join request for remote node %d at %s", nodeId, address)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel()
	err := node.SyncRequestAddNode(ctx, clusterId, nodeId, address, 0)
	if err != nil {
		raftLogger.Printf("join error: %d, %s", nodeId, address)
	}
	return err
}

func HandleDelete(nodeId uint64) error {
	raftLogger.Printf("received delete request node %d", nodeId)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel()
	err := node.SyncRequestDeleteNode(ctx, clusterId, nodeId, 0)
	if err != nil {
		raftLogger.Printf("delete error: %d, %v", nodeId, err)
	}
	return err
}

func GetNodes() ([]*pb.Node, error) {
	leader, _, err := node.GetLeaderID(clusterId)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.Conf.Cluster.Timeout)*time.Millisecond)
	defer cancel()
	membership, err := node.SyncGetClusterMembership(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Node, 0)
	for nodeId, address := range membership.Nodes {
		res = append(res, &pb.Node{Id: nodeId, Address: address, Leader: leader == nodeId})
	}
	return res, nil
}
