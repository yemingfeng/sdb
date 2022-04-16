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
