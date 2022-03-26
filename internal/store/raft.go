package store

import (
	"errors"
	"fmt"
	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
	"github.com/yemingfeng/sdb/internal/conf"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

var clusterId = uint64(1)
var node *dragonboat.NodeHost

func StartRaft() {
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
	if len(conf.Conf.Cluster.Master) == 0 {
		initialMembers := map[uint64]string{
			conf.Conf.Cluster.NodeId: conf.Conf.Cluster.Address,
		}
		if err := node.StartOnDiskCluster(initialMembers, false, NewFSM, rc); err != nil {
			log.Fatalln(err)
		}
	} else {
		if err := node.StartOnDiskCluster(nil, true, NewFSM, rc); err != nil {
			log.Fatalln(err)
		}
		if conf.Conf.Cluster.Join {
			resp, err := http.Get(fmt.Sprintf("http://%s/join?nodeId=%d&address=%s", conf.Conf.Cluster.Master, conf.Conf.Cluster.NodeId, conf.Conf.Cluster.Address))
			if err != nil {
				log.Fatalln(err)
			}
			bs, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			if "ok" != string(bs) {
				log.Fatalln("join failed")
			}
			defer func() {
				_ = resp.Body.Close()
			}()
		}
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

func HandleJoin(nodeId uint64, address string) error {
	log.Printf("received join request for remote node %d at %s", nodeId, address)
	rs, err := node.RequestAddNode(clusterId, nodeId, address, 0, time.Duration(conf.Conf.Cluster.Timeout)*time.Millisecond)
	if err != nil {
		log.Printf("join error: %d, %s", nodeId, address)
		return err
	}
	select {
	case r := <-rs.AppliedC():
		if r.Completed() {
			fmt.Printf("membership change completed successfully")
			return nil
		} else {
			return errors.New("join failed")
		}
	}
}
