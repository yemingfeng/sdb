package test

import (
	"fmt"
	"github.com/yemingfeng/sdb/internal/collection/list"
	"github.com/yemingfeng/sdb/internal/config"
	"github.com/yemingfeng/sdb/internal/service"
	"github.com/yemingfeng/sdb/internal/store"
	"testing"
)

func TestList(t *testing.T) {
	config := config.NewTestConfig()
	list := list.NewList()
	locker := service.NewLocker(config)
	db := store.NewStore(config)
	listService := service.NewListService(db, locker, list)

	for i := 0; i < 2; i++ {
		userKey := []byte(fmt.Sprintf("userKey:%d", i))
		existRes, err := listService.Exist(userKey)
		if err != nil {
			t.Fatalf("err: %+v", err)
		}
		t.Logf("existRes: %+v", existRes)

		err = listService.LPush(userKey, [][]byte{[]byte("v1"), []byte("v2"), []byte("v2"), []byte("v3")}, true)
		if err != nil {
			t.Fatalf("err: %+v", err)
		}

		err = listService.LPush(userKey, [][]byte{[]byte("v-1")}, true)
		if err != nil {
			t.Fatalf("err: %+v", err)
		}

		err = listService.RPush(userKey, [][]byte{[]byte("v4"), []byte("v5")}, true)
		if err != nil {
			t.Fatalf("err: %+v", err)
		}

		count, err := listService.Count(userKey)
		if err != nil {
			t.Fatalf("err: %+v", err)
		}
		t.Logf("Count: %d", count)

		rangeRes, err := listService.Range(userKey, 0, 10)
		if err != nil {
			t.Fatalf("err: %+v", err)
		}
		t.Logf("rangeRes: %+v", rangeResToPrintableString(rangeRes))

		rangeRes, err = listService.Range(userKey, 3, 2)
		if err != nil {
			t.Fatalf("err: %+v", err)
		}
		t.Logf("rangeRes: %+v", rangeResToPrintableString(rangeRes))

		rangeRes, err = listService.Range(userKey, -1, 5)
		if err != nil {
			t.Fatalf("err: %+v", err)
		}
		t.Logf("rangeRes: %+v", rangeResToPrintableString(rangeRes))

		getRes, err := listService.Get(userKey, []int64{0, 1, 2, 3, 4, 5})
		if err != nil {
			t.Fatalf("err: %+v", err)
		}
		t.Logf("getRes: %+v", getResToPrintableString(getRes))

		containRes, err := listService.Contain(userKey, [][]byte{[]byte("v1"), []byte("v2"), []byte("v3"), []byte("10")})
		if err != nil {
			t.Fatalf("err: %+v", err)
		}
		t.Logf("containRes: %+v", containRes)

		existRes, err = listService.Exist(userKey)
		if err != nil {
			t.Fatalf("err: %+v", err)
		}
		t.Logf("existRes: %+v", existRes)

		err = listService.Rem(userKey, []byte("v2"), true)
		if err != nil {
			t.Fatalf("err: %+v", err)
		}

		rangeRes, err = listService.Range(userKey, 0, 10)
		if err != nil {
			t.Fatalf("err: %+v", err)
		}
		t.Logf("rangeRes: %+v", rangeResToPrintableString(rangeRes))

		err = listService.Delete(userKey, true)
		if err != nil {
			t.Fatalf("err: %+v", err)
		}

		rangeRes, err = listService.Range(userKey, 0, 10)
		if err != nil {
			t.Fatalf("err: %+v", err)
		}
		t.Logf("rangeRes: %+v", rangeResToPrintableString(rangeRes))

		t.Log("==================")
	}
}

func getResToPrintableString(getRes map[int64][]byte) string {
	str := ""
	for index, value := range getRes {
		str += fmt.Sprintf("%d=%s, ", index, value)
	}
	return str
}

func rangeResToPrintableString(rangeRes [][]byte) string {
	str := ""
	for _, item := range rangeRes {
		str += string(item) + ", "
	}
	return str
}
