package service

import (
	"github.com/cockroachdb/pebble"
	"github.com/yemingfeng/sdb/internal/collection/list"
)

type ListService struct {
	db     *pebble.DB
	locker *Locker
	list   *list.List
}

func NewListService(db *pebble.DB, lock *Locker, list *list.List) *ListService {
	return &ListService{db: db, locker: lock, list: list}
}

func (listService *ListService) LPush(userKey []byte, values [][]byte, sync bool) error {
	listService.locker.wLock(userKey)
	defer listService.locker.wUnLock(userKey)

	batch := listService.db.NewIndexedBatch()
	defer batch.Close()

	if err := listService.list.LPush(batch, userKey, values); err != nil {
		return err
	}
	return batch.Commit(&pebble.WriteOptions{Sync: sync})
}

func (listService *ListService) RPush(userKey []byte, values [][]byte, sync bool) error {
	listService.locker.wLock(userKey)
	defer listService.locker.wUnLock(userKey)

	batch := listService.db.NewIndexedBatch()
	defer batch.Close()

	if err := listService.list.RPush(batch, userKey, values); err != nil {
		return err
	}
	return batch.Commit(&pebble.WriteOptions{Sync: sync})
}

func (listService *ListService) LPop(userKey []byte, sync bool) ([]byte, error) {
	listService.locker.wLock(userKey)
	defer listService.locker.wUnLock(userKey)

	batch := listService.db.NewIndexedBatch()
	defer batch.Close()

	if res, err := listService.list.LPop(batch, userKey); err != nil {
		return nil, err
	} else {
		if err := batch.Commit(&pebble.WriteOptions{Sync: sync}); err != nil {
			return nil, err
		}
		return res, nil
	}
}

func (listService *ListService) Rpop(userKey []byte, sync bool) ([]byte, error) {
	listService.locker.wLock(userKey)
	defer listService.locker.wUnLock(userKey)

	batch := listService.db.NewIndexedBatch()
	defer batch.Close()

	if res, err := listService.list.RPop(batch, userKey); err != nil {
		return nil, err
	} else {
		if err := batch.Commit(&pebble.WriteOptions{Sync: sync}); err != nil {
			return nil, err
		}
		return res, nil
	}
}

func (listService *ListService) Count(userKey []byte) (uint64, error) {
	listService.locker.rLock(userKey)
	defer listService.locker.rUnLock(userKey)

	batch := listService.db.NewIndexedBatch()
	defer batch.Close()

	if res, err := listService.list.Count(batch, userKey); err != nil {
		return 0, err
	} else {
		return res, nil
	}
}

func (listService *ListService) Range(userKey []byte, offset int64, limit uint64) ([][]byte, error) {
	listService.locker.rLock(userKey)
	defer listService.locker.rUnLock(userKey)

	batch := listService.db.NewIndexedBatch()
	defer batch.Close()

	if res, err := listService.list.Range(batch, userKey, offset, limit); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (listService *ListService) Get(userKey []byte, indexes []int64) (map[int64][]byte, error) {
	listService.locker.rLock(userKey)
	defer listService.locker.rUnLock(userKey)

	batch := listService.db.NewIndexedBatch()
	defer batch.Close()

	if res, err := listService.list.Get(batch, userKey, indexes); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (listService *ListService) Contain(userKey []byte, values [][]byte) ([]bool, error) {
	listService.locker.rLock(userKey)
	defer listService.locker.rUnLock(userKey)

	batch := listService.db.NewIndexedBatch()
	defer batch.Close()

	if res, err := listService.list.Contain(batch, userKey, values); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

func (listService *ListService) Exist(userKey []byte) (bool, error) {
	listService.locker.rLock(userKey)
	defer listService.locker.rUnLock(userKey)

	batch := listService.db.NewIndexedBatch()
	defer batch.Close()

	if res, err := listService.list.Exist(batch, userKey); err != nil {
		return false, err
	} else {
		return res, nil
	}
}

func (listService *ListService) Delete(userKey []byte, sync bool) error {
	listService.locker.wLock(userKey)
	defer listService.locker.wUnLock(userKey)

	batch := listService.db.NewIndexedBatch()
	defer batch.Close()

	if err := listService.list.Delete(batch, userKey); err != nil {
		return err
	} else {
		return batch.Commit(&pebble.WriteOptions{Sync: sync})
	}
}

func (listService *ListService) Rem(userKey []byte, value []byte, sync bool) error {
	listService.locker.wLock(userKey)
	defer listService.locker.wUnLock(userKey)

	batch := listService.db.NewIndexedBatch()
	defer batch.Close()

	if err := listService.list.Rem(batch, userKey, value); err != nil {
		return err
	} else {
		return batch.Commit(&pebble.WriteOptions{Sync: sync})
	}
}

func (listService *ListService) LInsert(userKey []byte, insertValue []byte, index int64, sync bool) error {
	listService.locker.wLock(userKey)
	defer listService.locker.wUnLock(userKey)

	batch := listService.db.NewIndexedBatch()
	defer batch.Close()

	if err := listService.list.LInsert(batch, userKey, insertValue, index); err != nil {
		return err
	} else {
		return batch.Commit(&pebble.WriteOptions{Sync: sync})
	}
}

func (listService *ListService) RInsert(userKey []byte, insertValue []byte, index int64, sync bool) error {
	listService.locker.wLock(userKey)
	defer listService.locker.wUnLock(userKey)

	batch := listService.db.NewIndexedBatch()
	defer batch.Close()

	if err := listService.list.RInsert(batch, userKey, insertValue, index); err != nil {
		return err
	} else {
		return batch.Commit(&pebble.WriteOptions{Sync: sync})
	}
}
