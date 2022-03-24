package store

import (
	"encoding/binary"
	"github.com/hashicorp/raft"
	"github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
)

var lastApplyIndexKey = []byte("last_apply_index_key")

type FSM struct {
	lastApplyIndex uint64
}

func NewFSM() *FSM {
	batch := NewBatch()
	defer batch.Close()

	lastApplyIndexValue, err := batch.Get(lastApplyIndexKey)
	if err != nil {
		log.Fatalf("get last apply index key error, %+v", err)
	}

	if len(lastApplyIndexValue) > 0 {
		lastApplyIndex := binary.LittleEndian.Uint64(lastApplyIndexValue)
		log.Printf("last apply index: %d", lastApplyIndex)
		return &FSM{lastApplyIndex: lastApplyIndex}
	}
	log.Println("not found last apply index")
	return &FSM{lastApplyIndex: 0}
}

func (fsm *FSM) Apply(raftLog *raft.Log) interface{} {
	pbLog := &pb.Log{}
	err := proto.Unmarshal(raftLog.Data, pbLog)
	if err != nil {
		return err
	}
	if len(pbLog.LogEntries) == 0 {
		return nil
	}

	batch := NewBatch()
	defer batch.Close()

	for _, logEntry := range pbLog.LogEntries {
		switch logEntry.Op {
		case pb.Op_OP_SET:
			if err := batch.Set(logEntry.Key, logEntry.Value); err != nil {
				return err
			}
			break
		case pb.Op_OP_DEL:
			if err := batch.Del(logEntry.Key); err != nil {
				return err
			}
			break
		}
	}

	if err := batch.Set(lastApplyIndexKey, util.UInt64ToBytes(raftLog.Index)); err != nil {
		return err
	}

	return batch.ApplyCommit()
}

func (fsm *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &Snapshot{lastApplyIndex: fsm.lastApplyIndex}, nil
}

func (fsm *FSM) Restore(_ io.ReadCloser) error {
	log.Println("fsm restore")
	return nil
}

type Snapshot struct {
	lastApplyIndex uint64
}

func (snapshot Snapshot) Persist(sink raft.SnapshotSink) error {
	log.Println("snapshot persist")
	defer log.Println("snapshot persist successful")

	_, err := sink.Write(util.UInt64ToBytes(snapshot.lastApplyIndex))
	return err
}

func (snapshot Snapshot) Release() {
	log.Println("snapshot release")
}
