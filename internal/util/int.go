package util

import (
	"encoding/binary"
)

func ToBytes(s int32) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(s))
	return bs
}

func UInt32ToBytes(s uint32) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, s)
	return bs
}

func UInt64ToBytes(s uint64) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, s)
	return bs
}

func BytesToUInt64(bs []byte) uint64 {
	return binary.LittleEndian.Uint64(bs)
}
