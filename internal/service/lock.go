package service

import (
	"github.com/howeyc/crc16"
	"sync"
)

var lockers = make(map[DataType][]sync.RWMutex)

type DataType int8

const (
	LString      DataType = 0x0
	LList        DataType = 0x1
	LSet         DataType = 0x2
	LSortedSet   DataType = 0x3
	LBloomFilter DataType = 0x4
	LHyperLogLog DataType = 0x5
	LBitset      DataType = 0x6
	LMap         DataType = 0x7
	LGeoHash     DataType = 0x8
)

var lockerCounts = map[DataType]int{
	LString:      16,
	LList:        16,
	LSet:         16,
	LSortedSet:   16,
	LBloomFilter: 16,
	LHyperLogLog: 16,
	LBitset:      16,
	LMap:         16,
	LGeoHash:     16,
}

func init() {
	for dataType, lockerCount := range lockerCounts {
		lockers[dataType] = make([]sync.RWMutex, lockerCount)
		for i := 0; i < lockerCount; i++ {
			lockers[dataType][i] = sync.RWMutex{}
		}
	}
}

func lock(dataType DataType, key []byte) {
	locker := getLocker(dataType, key)
	locker.Lock()
}

func unlock(dataType DataType, key []byte) {
	locker := getLocker(dataType, key)
	locker.Unlock()
}

func getLocker(dataType DataType, key []byte) *sync.RWMutex {
	checksum := crc16.Checksum(key, crc16.IBMTable)
	return &lockers[dataType][int(checksum)%lockerCounts[dataType]]
}
