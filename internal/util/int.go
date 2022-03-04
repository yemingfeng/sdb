package util

import (
	"encoding/binary"
	"strconv"
)

func ToBytes(s int32) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(s))
	return bs
}

func StringToInt32(s string) (int32, error) {
	t, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return int32(t), err
}

func StringToUInt32(s string) (uint32, error) {
	t, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint32(t), err
}

func StringToDouble(s string) (float64, error) {
	t, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return t, err
}

func StringToBoolean(s string) (bool, error) {
	t, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return false, err
	}
	if t > 0 {
		return true, err
	}
	return false, err
}
