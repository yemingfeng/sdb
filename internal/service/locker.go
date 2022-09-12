package service

import (
	"github.com/yemingfeng/sdb/internal/config"
	"hash/fnv"
	"sync"
)

type Locker struct {
	lockers []*sync.RWMutex
}

func NewLocker(config *config.Config) *Locker {
	lockers := make([]*sync.RWMutex, config.Collection.LockerCount)
	for i := 0; i < len(lockers); i++ {
		lockers[i] = &sync.RWMutex{}
	}
	return &Locker{
		lockers: lockers,
	}
}

func (locker *Locker) batchLock(userKeys [][]byte) {
	hashes := make([]bool, len(locker.lockers))
	for i := 0; i < len(userKeys); i++ {
		hashes[locker.hash(userKeys[i])] = true
	}
	for i := 0; i < len(hashes); i++ {
		if hashes[i] {
			locker.lockers[i].Lock()
		}
	}
}

func (locker *Locker) batchUnLock(userKeys [][]byte) {
	hashes := make([]bool, len(locker.lockers))
	for i := 0; i < len(userKeys); i++ {
		hashes[locker.hash(userKeys[i])] = true
	}
	for i := 0; i < len(hashes); i++ {
		if hashes[i] {
			locker.lockers[i].Unlock()
		}
	}
}

func (locker *Locker) hash(userKey []byte) int {
	h := fnv.New32a()
	_, _ = h.Write(userKey)
	return int(h.Sum32()) % len(locker.lockers)
}

func (locker *Locker) getLocker(userKey []byte) *sync.RWMutex {
	return locker.lockers[locker.hash(userKey)]
}

func (locker *Locker) rLock(userKey []byte) {
	locker.getLocker(userKey).RLock()
}

func (locker *Locker) rUnLock(userKey []byte) {
	locker.getLocker(userKey).RUnlock()
}

func (locker *Locker) wLock(userKey []byte) {
	locker.getLocker(userKey).Lock()
}

func (locker *Locker) wUnLock(userKey []byte) {
	locker.getLocker(userKey).Unlock()
}
