package store

import (
	"github.com/lni/dragonboat/v3/statemachine"
	"github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"google.golang.org/protobuf/proto"
	"io"
)

var fsmLogger = util.GetLogger("fsm")
var lastApplyIndexKey = []byte("last_apply_index_key")

type FSM struct {
	lastApplyIndex uint64
}

func NewFSM(_ uint64, _ uint64) statemachine.IOnDiskStateMachine {
	return &FSM{}
}

func (fsm FSM) Open(stopc <-chan struct{}) (uint64, error) {
	batch := NewBatch()
	defer batch.Close()

	lastApplyIndexValue, err := batch.Get(lastApplyIndexKey)
	if err != nil {
		fsmLogger.Fatalf("get last apply index key error, %+v", err)
		return 0, err
	}

	if len(lastApplyIndexValue) > 0 {
		fsm.lastApplyIndex = util.BytesToUInt64(lastApplyIndexValue)
		fsmLogger.Printf("last apply index: %d", fsm.lastApplyIndex)
		return fsm.lastApplyIndex, nil
	}
	fsmLogger.Println("not found last apply index")
	fsm.lastApplyIndex = uint64(0)
	return fsm.lastApplyIndex, nil
}

func (fsm FSM) Update(entries []statemachine.Entry) ([]statemachine.Entry, error) {
	batch := NewBatch()
	defer batch.Close()

	for _, entry := range entries {
		pbLog := &pb.Log{}
		err := proto.Unmarshal(entry.Cmd, pbLog)
		if err != nil {
			fsmLogger.Printf("error on unmarshal pbLog: [%+v], err: [%+v]", pbLog, err)
			return entries, err
		}
		for _, logEntry := range pbLog.LogEntries {
			switch logEntry.Op {
			case pb.Op_OP_SET:
				if err := batch.Set(logEntry.Key, logEntry.Value); err != nil {
					return entries, err
				}
				break
			case pb.Op_OP_DEL:
				if err := batch.Del(logEntry.Key); err != nil {
					return entries, err
				}
				break
			}
		}
		entry.Result = statemachine.Result{Value: uint64(len(entry.Cmd))}
	}

	if err := batch.Set(lastApplyIndexKey, util.UInt64ToBytes(entries[len(entries)-1].Index)); err != nil {
		return entries, err
	}

	return entries, batch.ApplyCommit()
}

func (fsm FSM) Lookup(key interface{}) (interface{}, error) {
	batch := NewBatch()
	defer batch.Close()
	return batch.Get(key.([]byte))
}

func (fsm FSM) Sync() error {
	return nil
}

func (fsm FSM) PrepareSnapshot() (interface{}, error) {
	return fsm.lastApplyIndex, nil
}

func (fsm FSM) SaveSnapshot(i interface{}, writer io.Writer, i2 <-chan struct{}) error {
	_, err := writer.Write(util.UInt64ToBytes(i.(uint64)))
	fsmLogger.Printf("save snapshot, last apply index: [%d]", fsm.lastApplyIndex)
	return err
}

func (fsm FSM) RecoverFromSnapshot(reader io.Reader, i <-chan struct{}) error {
	bs := make([]byte, 8)
	_, err := reader.Read(bs)
	if err != nil {
		return err
	}
	fsm.lastApplyIndex = util.BytesToUInt64(bs)
	fsmLogger.Printf("recover from snapshot, last apply index: [%d]", fsm.lastApplyIndex)
	return nil
}

func (fsm FSM) Close() error {
	batch := NewBatch()
	defer batch.Close()
	fsmLogger.Printf("close, last apply index: [%d]", fsm.lastApplyIndex)
	return batch.Set(lastApplyIndexKey, util.UInt64ToBytes(fsm.lastApplyIndex))
}
