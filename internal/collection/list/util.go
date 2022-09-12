package list

import (
	"errors"
	"fmt"
	"github.com/yemingfeng/sdb/internal/util"
	"math"
	"strconv"
	"strings"
)

var listLogger = util.GetLogger("list")

const metaKeyMagic = "am"
const seqKeyMagic = "as"
const valueKeyMagic = "av"

// generateMetaKey, return am:{userKey}:
func generateMetaKey(userKey []byte, version uint64) []byte {
	return []byte(fmt.Sprintf("%s:%s:%064d", metaKeyMagic, util.Base64Encode(userKey), version))
}

// generateMetaPrefixKey, return am:{userKey}:
func generateMetaPrefixKey(userKey []byte) []byte {
	return []byte(fmt.Sprintf("%s:%s:", metaKeyMagic, util.Base64Encode(userKey)))
}

// parseMetaKey, return userKey, version, error
func parseMetaKey(metaKey []byte) ([]byte, uint64, error) {
	infos := strings.Split(string(metaKey), ":")
	if len(infos) != 3 {
		listLogger.Printf("can not parse metaKey: %s, length != 3", metaKey)
		return nil, 0, errors.New(fmt.Sprintf("can not parse metaKey: %s", metaKey))
	}
	if infos[0] != metaKeyMagic {
		listLogger.Printf("can not parse metaKey: %s, %s is not %s", metaKey, infos[0], metaKeyMagic)
		return nil, 0, errors.New(fmt.Sprintf("can not parse metaKey: %s, %s is not %s", metaKey, infos[0], metaKeyMagic))
	}
	userKey, err := util.Base64Decode([]byte(infos[1]))
	if err != nil {
		listLogger.Printf("can not base64Encode metaKey on userKey: %s, error: %+v", metaKey, err)
		return nil, 0, err
	}
	version, err := strconv.ParseUint(infos[2], 10, 64)
	if err != nil {
		listLogger.Printf("can not parseUint metaKey on seq: %s, error: %+v", metaKey, err)
		return nil, 0, err
	}
	return userKey, version, nil
}

// generateSeqKey, return as:{userKey}:{version}:{seq}:{value}
func generateSeqKey(userKey []byte, version uint64, seq int64, value []byte) []byte {
	return []byte(fmt.Sprintf("%s:%s:%064d:%064d:%s", seqKeyMagic, util.Base64Encode(userKey), version, seq+math.MaxInt32, util.Base64Encode(value)))
}

// generateSeqPrefixKey, return as:{userKey}:{version}:{seq}:
func generateSeqPrefixKey(userKey []byte, version uint64, seq int64) []byte {
	return []byte(fmt.Sprintf("%s:%s:%064d:%064d:", seqKeyMagic, util.Base64Encode(userKey), version, seq+math.MaxInt32))
}

// generateSeqIteratorKey, return as:{userKey}:{version}:
func generateSeqIteratorKey(userKey []byte, version uint64) []byte {
	return []byte(fmt.Sprintf("%s:%s:%064d:", seqKeyMagic, util.Base64Encode(userKey), version))
}

// parseSeqKey, return userKey, version, seq, value, error
func parseSeqKey(seqKey []byte) ([]byte, uint64, int64, []byte, error) {
	infos := strings.Split(string(seqKey), ":")
	if len(infos) != 5 {
		listLogger.Printf("can not parse seqKey: %s, length != 5", seqKey)
		return nil, 0, 0, nil, errors.New(fmt.Sprintf("can not parse seqKey: %s, length != 5", seqKey))
	}
	if infos[0] != seqKeyMagic {
		listLogger.Printf("can not parse seqKey: %s, %s is not %s", seqKey, infos[0], seqKeyMagic)
		return nil, 0, 0, nil, errors.New(fmt.Sprintf("can not parse seqKey: %s, %s is not %s", seqKey, infos[0], seqKeyMagic))
	}
	userKey, err := util.Base64Decode([]byte(infos[1]))
	if err != nil {
		listLogger.Printf("can not base64Encode seqKey on userKey: %s, error: %+v", seqKey, err)
		return nil, 0, 0, nil, err
	}
	version, err := strconv.ParseUint(infos[2], 10, 64)
	if err != nil {
		listLogger.Printf("can not parseUint seqKey on version: %s, error: %+v", seqKey, err)
		return nil, 0, 0, nil, err
	}
	seq, err := strconv.ParseInt(infos[3], 10, 64)
	if err != nil {
		listLogger.Printf("can not parseInt seqKey on seq: %s, error: %+v", seqKey, err)
		return nil, 0, 0, nil, err
	}
	value, err := util.Base64Decode([]byte(infos[4]))
	if err != nil {
		listLogger.Printf("can not base64Encode seqKey on value: %s, error: %+v", seqKey, err)
		return nil, 0, 0, nil, err
	}
	return userKey, version, seq - math.MaxInt32, value, nil
}

// generateValueKey, return av:{userKey}:{version}:{value}:{seq}
func generateValueKey(userKey []byte, version uint64, value []byte, seq int64) []byte {
	return []byte(fmt.Sprintf("%s:%s:%064d:%s:%064d", valueKeyMagic, util.Base64Encode(userKey), version, util.Base64Encode(value), seq+math.MaxInt32))
}

// generateValuePrefixKey, return av:{userKey}:{version}:{value}:
func generateValuePrefixKey(userKey []byte, version uint64, value []byte) []byte {
	return []byte(fmt.Sprintf("%s:%s:%064d:%s:", valueKeyMagic, util.Base64Encode(userKey), version, util.Base64Encode(value)))
}

// parseValueKey, return userKey, version, value, seq, error
func parseValueKey(valueKey []byte) ([]byte, uint64, []byte, int64, error) {
	infos := strings.Split(string(valueKey), ":")
	if len(infos) != 5 {
		listLogger.Printf("can not parse valueKey: %s, length != 5", valueKey)
		return nil, 0, nil, 0, errors.New(fmt.Sprintf("can not parse valueKey: %s", valueKey))
	}
	if infos[0] != valueKeyMagic {
		listLogger.Printf("can not parse valueKey: %s, %s is not %s", valueKey, infos[0], valueKeyMagic)
		return nil, 0, nil, 0, errors.New(fmt.Sprintf("can not parse valueKey: %s, %s is not %s", valueKey, infos[0], valueKeyMagic))
	}
	userKey, err := util.Base64Decode([]byte(infos[1]))
	if err != nil {
		listLogger.Printf("can not base64Encode valueKey on userKey: %s, error: %+v", valueKey, err)
		return nil, 0, nil, 0, err
	}
	version, err := strconv.ParseUint(infos[2], 10, 64)
	if err != nil {
		listLogger.Printf("can not parseUint valueKey on seq: %s, error: %+v", valueKey, err)
		return nil, 0, nil, 0, err
	}
	value, err := util.Base64Decode([]byte(infos[3]))
	if err != nil {
		listLogger.Printf("can not base64Encode valueKey on value: %s, error: %+v", valueKey, err)
		return nil, 0, nil, 0, err
	}
	seq, err := strconv.ParseInt(infos[4], 10, 64)
	if err != nil {
		listLogger.Printf("can not parseInt valueKey on seq: %s, error: %+v", valueKey, err)
		return nil, 0, nil, 0, err
	}
	return userKey, version, value, seq - math.MaxInt32, nil
}
