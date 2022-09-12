package list

import (
	"errors"
	"fmt"
	"github.com/cockroachdb/pebble"
	"github.com/yemingfeng/sdb/internal/model"
	"github.com/yemingfeng/sdb/internal/store"
	"github.com/yemingfeng/sdb/internal/util"
	"google.golang.org/protobuf/proto"
	"reflect"
)

type List struct {
}

func NewList() *List {
	return &List{}
}

// LPush
// left push elements, such as if the list contains [a, b, c],
// now want to left push [d, e],
// then the list will update be: [e, d, a, b, c]
func (list *List) LPush(batch *pebble.Batch, userKey []byte, values [][]byte) error {
	return list.push(batch, userKey, values, true)
}

// RPush
// right push elements, such as if the list contains [a, b, c],
// now want to left push [d, e],
// then the list will update be: [a, b, c, d, e]
func (list *List) RPush(batch *pebble.Batch, userKey []byte, values [][]byte) error {
	return list.push(batch, userKey, values, false)
}

func (list *List) push(batch *pebble.Batch, userKey []byte, values [][]byte, left bool) error {
	if len(userKey) == 0 {
		return errors.New("userKey is emtpy")
	}
	meta, version, err := list.getMeta(batch, userKey)
	if err != nil {
		return err
	}
	if meta == nil {
		version = util.NextId()
		meta = &model.ListMeta{
			Head:    0,
			Tail:    1,
			Count:   0,
			Deleted: false,
		}
	}

	for i := 0; i < len(values); i++ {
		var seq int64
		if left {
			seq = meta.Head
			meta.Head--
		} else {
			seq = meta.Tail
			meta.Tail++
		}
		value := values[i]
		seqKey := generateSeqKey(userKey, version, seq, value)
		if err := batch.Set(seqKey, nil, nil); err != nil {
			return err
		}
		valueKey := generateValueKey(userKey, version, value, seq)
		if err := batch.Set(valueKey, nil, nil); err != nil {
			return err
		}
	}

	meta.Count = uint64(len(values)) + meta.Count

	metaValue, err := proto.Marshal(meta)
	if err != nil {
		return err
	}
	if err := batch.Set(generateMetaKey(userKey, version), metaValue, nil); err != nil {
		return err
	}
	return nil
}

func (list *List) LPop(batch *pebble.Batch, userKey []byte) ([]byte, error) {
	return list.pop(batch, userKey, true)
}

func (list *List) RPop(batch *pebble.Batch, userKey []byte) ([]byte, error) {
	return list.pop(batch, userKey, false)
}

func (list *List) pop(batch *pebble.Batch, userKey []byte, left bool) ([]byte, error) {
	if len(userKey) == 0 {
		return nil, errors.New("userKey is emtpy")
	}
	meta, version, err := list.getMeta(batch, userKey)
	if err != nil || meta == nil || meta.Count == 0 {
		return nil, err
	}

	seq := int64(0)
	if left {
		seq = meta.Head + 1
		meta.Head += 1
	} else {
		seq = meta.Tail - 1
		meta.Tail -= 1
	}

	seqKey := list.getSeqKey(batch, userKey, version, seq)
	_, _, _, value, err := parseSeqKey(seqKey)
	if err != nil {
		return nil, err
	}
	if err := batch.Delete(generateValueKey(userKey, version, value, seq), nil); err != nil {
		return nil, err
	}
	if err := batch.Delete(generateSeqKey(userKey, version, seq, value), nil); err != nil {
		return nil, err
	}
	meta.Count -= 1

	metaValue, err := proto.Marshal(meta)
	if err != nil {
		return nil, err
	}
	if err := batch.Set(generateMetaKey(userKey, version), metaValue, nil); err != nil {
		return nil, err
	}

	return value, nil
}

// Count
// return element counts, such as if the list contains [a, b, c], will return 3
func (list *List) Count(batch *pebble.Batch, userKey []byte) (uint64, error) {
	if len(userKey) == 0 {
		return 0, errors.New("userKey is emtpy")
	}
	meta, _, err := list.getMeta(batch, userKey)
	if err != nil || meta == nil {
		return 0, err
	}
	return meta.Count, nil
}

// Range
// range list, such as if the list contains [a, b, c, d, e]
// if pass offset=0,limit=3, will return [a, b, c]
// if pass offset=1,limit=2, will return [a, b]
// if pass offset=-1,limit=2, will return [e, d]
// if pass offset=-2,limit=5, will return [c, b, a]
func (list *List) Range(batch *pebble.Batch, userKey []byte, offset int64, limit uint64) ([][]byte, error) {
	if len(userKey) == 0 {
		return nil, errors.New("userKey is emtpy")
	}
	iter := batch.NewIter(&pebble.IterOptions{})
	defer iter.Close()

	meta, version, err := list.getMeta(batch, userKey)
	if err != nil || meta == nil || meta.Count == 0 {
		return nil, err
	}

	reverse := offset < 0
	if !reverse {
		lowerBound := generateSeqPrefixKey(userKey, version, meta.Head+offset+1)
		upperBound := generateSeqPrefixKey(userKey, version, meta.Head+offset+int64(limit)+1)
		iter.SetBounds(lowerBound, upperBound)
		iter.First()
	} else {
		lowerBound := generateSeqPrefixKey(userKey, version, meta.Tail+offset-int64(limit)+1)
		upperBound := generateSeqPrefixKey(userKey, version, meta.Tail+offset+1)
		iter.SetBounds(lowerBound, upperBound)
		iter.Last()
	}

	res := make([][]byte, limit)
	i := uint64(0)

	for {
		if !iter.Valid() {
			break
		}
		seqKey := iter.Key()
		_, _, _, value, err := parseSeqKey(seqKey)

		if err != nil {
			return nil, err
		}
		res[i] = value
		i++
		if !reverse {
			iter.Next()
		} else {
			iter.Prev()
		}
	}
	return res[:i], nil
}

// Get
// return get elements, such as if the list contains [a, b, c], if pass index = [1, 2], return [b, c]
func (list *List) Get(batch *pebble.Batch, userKey []byte, indexes []int64) (map[int64][]byte, error) {
	if len(userKey) == 0 {
		return nil, errors.New("userKey is emtpy")
	}
	meta, version, err := list.getMeta(batch, userKey)
	if err != nil || meta == nil || meta.Count == 0 {
		return nil, err
	}

	res := make(map[int64][]byte)
	for i := 0; i < len(indexes); i++ {
		index := indexes[i]
		seq := meta.Head + 1 + index
		if seq > meta.Tail-1 || seq < meta.Head {
			return nil, errors.New(fmt.Sprintf("invalid index: %d", index))
		}
		seqKey := list.getSeqKey(batch, userKey, version, seq)
		if len(seqKey) > 0 {
			_, _, _, value, err := parseSeqKey(seqKey)
			if err != nil {
				return nil, err
			}
			res[index] = value
		}
	}

	return res, nil
}

// Contain
// return dose elements if not contains, such as if the list contains [a, b, c], if pass index = [a, b, d], return [true, true, false]
func (list *List) Contain(batch *pebble.Batch, userKey []byte, values [][]byte) ([]bool, error) {
	if len(userKey) == 0 {
		return nil, errors.New("userKey is emtpy")
	}
	iter := batch.NewIter(&pebble.IterOptions{})
	defer iter.Close()

	meta, version, err := list.getMeta(batch, userKey)
	if err != nil || meta == nil || meta.Count == 0 {
		return nil, err
	}

	res := make([]bool, len(values))
	for i := 0; i < len(values); i++ {
		value := values[i]
		iterOptions := store.NewPrefixIterOptions(generateValuePrefixKey(userKey, version, value))
		iter.SetBounds(iterOptions.LowerBound, iterOptions.UpperBound)
		iter.First()
		res[i] = iter.Valid()
	}

	return res, nil
}

// Exist
// return dose elements if not exist
func (list *List) Exist(batch *pebble.Batch, userKey []byte) (bool, error) {
	if len(userKey) == 0 {
		return false, errors.New("userKey is emtpy")
	}
	iter := batch.NewIter(&pebble.IterOptions{})
	defer iter.Close()

	meta, _, err := list.getMeta(batch, userKey)
	if err != nil || meta == nil {
		return false, err
	}
	return true, nil
}

// Delete
// delete list
func (list *List) Delete(batch *pebble.Batch, userKey []byte) error {
	if len(userKey) == 0 {
		return errors.New("userKey is emtpy")
	}
	meta, version, err := list.getMeta(batch, userKey)
	if err != nil || meta == nil {
		return err
	}

	meta.Deleted = true

	metaValue, err := proto.Marshal(meta)
	if err != nil {
		return err
	}
	if err := batch.Set(generateMetaKey(userKey, version), metaValue, nil); err != nil {
		return err
	}

	return nil
}

// Rem
// if list contains [a, b, c, d, a, e], then pass 'a'
// the list will be [b, c, d, e]
func (list *List) Rem(batch *pebble.Batch, userKey []byte, value []byte) error {
	if len(userKey) == 0 {
		return errors.New("userKey is emtpy")
	}
	meta, version, err := list.getMeta(batch, userKey)
	if err != nil || meta == nil || meta.Count == 0 {
		return err
	}

	iter := batch.NewIter(store.NewPrefixIterOptions(generateValuePrefixKey(userKey, version, value)))
	defer iter.Close()

	startSeq := int64(0)
	found := false
	iter.First()
	if iter.Valid() {
		valueKey := iter.Key()
		_, _, _, seq, err := parseValueKey(valueKey)
		if err != nil {
			return err
		}
		startSeq = seq
		found = true
	}
	if !found {
		return nil
	}

	lowerBound := generateSeqPrefixKey(userKey, version, startSeq)
	upperBound := store.Next(generateSeqIteratorKey(userKey, version))
	iter.SetBounds(lowerBound, upperBound)
	foundCount := int64(0)
	//preSeq := startSeq
	// [a, a, a, b, a, a, c, d, a]
	// [1, 2, 3, 4, 5, 6, 7, 8, 9]
	for iter.First(); iter.Valid(); iter.Next() {
		// get current seq
		_, _, currentSeq, currentValue, err := parseSeqKey(iter.Key())
		if err != nil {
			return err
		}

		// delete old seq
		if err := batch.Delete(generateSeqKey(userKey, version, currentSeq, currentValue), nil); err != nil {
			return err
		}
		if err := batch.Delete(generateValueKey(userKey, version, currentValue, currentSeq), nil); err != nil {
			return err
		}
		if reflect.DeepEqual(currentValue, value) {
			foundCount++
		} else {
			if err := batch.Set(generateValueKey(userKey, version, currentValue, currentSeq-foundCount), nil, nil); err != nil {
				return err
			}
			if err := batch.Set(generateSeqKey(userKey, version, currentSeq-foundCount, currentValue), nil, nil); err != nil {
				return err
			}
		}
	}

	meta.Count = meta.Count - uint64(foundCount)
	meta.Tail = meta.Tail - foundCount

	metaValue, err := proto.Marshal(meta)
	if err != nil {
		return err
	}
	if err := batch.Set(generateMetaKey(userKey, version), metaValue, nil); err != nil {
		return err
	}

	return nil
}

func (list *List) LInsert(batch *pebble.Batch, userKey []byte, insertValue []byte, index int64) error {
	return list.insert(batch, userKey, insertValue, index, true)
}

func (list *List) RInsert(batch *pebble.Batch, userKey []byte, insertValue []byte, index int64) error {
	return list.insert(batch, userKey, insertValue, index, false)
}

// Insert
// if list contains [a, b, c, d, a, e], then pass insertValue='a', index=2
// if pass left=true, the list will be [a, b, a, c, d, a, e]
// if pass left=false, the list will be [a, b, c, a, d, a, e]
// if list not exist, will create new
func (list *List) insert(batch *pebble.Batch, userKey []byte, insertValue []byte, index int64, left bool) error {
	if len(userKey) == 0 {
		return errors.New("userKey is emtpy")
	}
	meta, version, err := list.getMeta(batch, userKey)
	if err != nil {
		return err
	}
	if meta == nil {
		version = util.NextId()
		meta = &model.ListMeta{
			Head:    0,
			Tail:    1,
			Count:   0,
			Deleted: false,
		}
	}

	// if index=0 & count=0, then just insert
	if index == 0 && meta.Count == 0 {
		seq := meta.Tail
		seqKey := generateSeqKey(userKey, version, seq, insertValue)
		if err := batch.Set(seqKey, nil, nil); err != nil {
			return err
		}
		valueKey := generateValueKey(userKey, version, insertValue, seq)
		if err := batch.Set(valueKey, nil, nil); err != nil {
			return err
		}
		meta.Tail++
		meta.Count++

		metaValue, err := proto.Marshal(meta)
		if err != nil {
			return err
		}
		if err := batch.Set(generateMetaKey(userKey, version), metaValue, nil); err != nil {
			return err
		}
		return nil
	}

	insertSeq := meta.Head + 1 + index
	if insertSeq < meta.Head || insertSeq > meta.Tail {
		return errors.New(fmt.Sprintf("invalid index: %d", index))
	}

	// move [meta.Tail...insertSeq+1]
	iter := batch.NewIter(&pebble.IterOptions{
		LowerBound: generateSeqPrefixKey(userKey, version, insertSeq+1),
		UpperBound: generateSeqPrefixKey(userKey, version, meta.Tail),
	})
	defer iter.Close()
	for iter.Last(); iter.Valid(); iter.Prev() {
		currentSeqKey := iter.Key()
		_, _, currentSeq, currentValue, err := parseSeqKey(currentSeqKey)
		if err != nil {
			return err
		}
		// delete old currentKey & valueKey
		if err := batch.Delete(currentSeqKey, nil); err != nil {
			return err
		}
		if err := batch.Delete(generateValueKey(userKey, version, currentValue, currentSeq), nil); err != nil {
			return err
		}

		// insert to new seq
		if err := batch.Set(generateSeqKey(userKey, version, currentSeq+1, currentValue), nil, nil); err != nil {
			return err
		}
		if err := batch.Set(generateValueKey(userKey, version, currentValue, currentSeq+1), nil, nil); err != nil {
			return err
		}
	}
	insertSeqKey := list.getSeqKey(batch, userKey, version, insertSeq)
	_, _, _, insertSeqValue, err := parseSeqKey(insertSeqKey)
	if err != nil {
		return err
	}
	if err := batch.Delete(insertSeqKey, nil); err != nil {
		return err
	}
	if err := batch.Delete(generateValueKey(userKey, version, insertSeqValue, insertSeq), nil); err != nil {
		return err
	}
	if left {
		if err := batch.Set(generateSeqKey(userKey, version, insertSeq, insertValue), nil, nil); err != nil {
			return err
		}
		if err := batch.Set(generateValueKey(userKey, version, insertValue, insertSeq), nil, nil); err != nil {
			return err
		}
		if err := batch.Set(generateSeqKey(userKey, version, insertSeq+1, insertSeqValue), nil, nil); err != nil {
			return err
		}
		if err := batch.Set(generateValueKey(userKey, version, insertSeqValue, insertSeq+1), nil, nil); err != nil {
			return err
		}
	} else {
		if err := batch.Set(generateSeqKey(userKey, version, insertSeq+1, insertValue), nil, nil); err != nil {
			return err
		}
		if err := batch.Set(generateValueKey(userKey, version, insertValue, insertSeq+1), nil, nil); err != nil {
			return err
		}
		if err := batch.Set(generateSeqKey(userKey, version, insertSeq, insertSeqValue), nil, nil); err != nil {
			return err
		}
		if err := batch.Set(generateValueKey(userKey, version, insertSeqValue, insertSeq), nil, nil); err != nil {
			return err
		}
	}

	meta.Count += 1
	meta.Tail += 1

	metaValue, err := proto.Marshal(meta)
	if err != nil {
		return err
	}
	return batch.Set(generateMetaKey(userKey, version), metaValue, nil)
}

// getMeta
// return max version && valid meta
// if max version meta is valid, return nil
func (list *List) getMeta(batch *pebble.Batch, userKey []byte) (*model.ListMeta, uint64, error) {
	iter := batch.NewIter(store.NewPrefixIterOptions(generateMetaPrefixKey(userKey)))
	defer iter.Close()

	if res := iter.Last(); !res {
		return nil, 0, nil
	}

	metaKey := iter.Key()
	metaValue := iter.Value()

	if len(metaKey) == 0 {
		return nil, 0, nil
	}
	_, version, err := parseMetaKey(metaKey)
	if err != nil {
		return nil, 0, err
	}
	meta := &model.ListMeta{}
	if err := proto.Unmarshal(metaValue, meta); err != nil {
		return nil, 0, err
	}
	valid := meta.Deleted != true
	if valid {
		return meta, version, nil
	}

	return nil, 0, nil
}

func (list *List) getSeqKey(batch *pebble.Batch, userKey []byte, version uint64, seq int64) []byte {
	iter := batch.NewIter(store.NewPrefixIterOptions(generateSeqPrefixKey(userKey, version, seq)))
	defer iter.Close()

	if res := iter.Last(); !res {
		return nil
	}

	return iter.Key()
}
