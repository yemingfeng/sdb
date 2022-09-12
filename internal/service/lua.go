package service

import (
	"fmt"
	"github.com/cockroachdb/pebble"
	"github.com/yemingfeng/sdb/internal/collection/list"
	"github.com/yemingfeng/sdb/internal/util"
	lua "github.com/yuin/gopher-lua"
)

var luaLogger = util.GetLogger("lua")

type LuaService struct {
	db     *pebble.DB
	locker *Locker
	list   *list.List
}

func NewLuaService(db *pebble.DB, locker *Locker, list *list.List) *LuaService {
	return &LuaService{db: db, locker: locker, list: list}
}

func (luaService *LuaService) Execute(userKeys [][]byte, script string) ([]string, error) {
	L := lua.NewState()
	defer L.Close()

	luaService.locker.batchLock(userKeys)
	defer luaService.locker.batchUnLock(userKeys)

	batch := luaService.db.NewIndexedBatch()
	defer batch.Close()

	found := func(uk string) bool {
		for _, userKey := range userKeys {
			if uk == string(userKey) {
				return true
			}
		}
		return false
	}

	L.SetGlobal("LLPush", L.NewFunction(func(L *lua.LState) int {
		top := L.GetTop()
		userKey := L.CheckString(1)
		if !found(userKey) {
			L.RaiseError(fmt.Sprintf("can not opt not lock userKey: %s", userKey))
		}

		values := make([][]byte, top-1)
		for i := 0; i < len(values); i++ {
			values[i] = []byte(L.CheckString(i + 2))
		}
		luaLogger.Printf("[LLPush] userKey: %s, values: %s", userKey, values)

		err := luaService.list.LPush(batch, []byte(userKey), values)
		if err != nil {
			L.RaiseError("%s", err)
		}
		return 0
	}))
	L.SetGlobal("LRPush", L.NewFunction(func(L *lua.LState) int {
		top := L.GetTop()
		userKey := L.CheckString(1)
		if !found(userKey) {
			L.RaiseError(fmt.Sprintf("can not opt not lock userKey: %s", userKey))
		}

		values := make([][]byte, top-1)
		for i := 0; i < len(values); i++ {
			values[i] = []byte(L.CheckString(i + 2))
		}
		luaLogger.Printf("[LRPush] userKey: %s, values: %s", userKey, values)

		err := luaService.list.RPush(batch, []byte(userKey), values)
		if err != nil {
			L.RaiseError("%s", err)
		}
		return 0
	}))
	L.SetGlobal("LLPop", L.NewFunction(func(L *lua.LState) int {
		userKey := L.CheckString(1)
		if !found(userKey) {
			L.RaiseError(fmt.Sprintf("can not opt not lock userKey: %s", userKey))
		}

		luaLogger.Printf("[LLPop] userKey: %s", userKey)

		res, err := luaService.list.LPop(batch, []byte(userKey))
		if err != nil {
			L.RaiseError("%s", err)
		}
		L.Push(lua.LString(res))
		return 1
	}))
	L.SetGlobal("LRPop", L.NewFunction(func(L *lua.LState) int {
		userKey := L.CheckString(1)
		if !found(userKey) {
			L.RaiseError(fmt.Sprintf("can not opt not lock userKey: %s", userKey))
		}
		luaLogger.Printf("[LRPop] userKey: %s", userKey)

		res, err := luaService.list.RPop(batch, []byte(userKey))
		if err != nil {
			L.RaiseError("%s", err)
		}
		L.Push(lua.LString(res))
		return 1
	}))
	L.SetGlobal("LCount", L.NewFunction(func(L *lua.LState) int {
		userKey := L.CheckString(1)
		luaLogger.Printf("[LCount] userKey: %s", userKey)

		res, err := luaService.list.Count(batch, []byte(userKey))
		if err != nil {
			L.RaiseError("%s", err)
		}
		L.Push(lua.LNumber(res))
		return 1
	}))
	L.SetGlobal("LRange", L.NewFunction(func(L *lua.LState) int {
		userKey := L.CheckString(1)
		offset := L.CheckInt64(2)
		limit := uint64(L.CheckInt64(3))
		luaLogger.Printf("[LRange] userKey: %s, offset: %d, limit: %d", userKey, offset, limit)

		res, err := luaService.list.Range(batch, []byte(userKey), offset, limit)
		if err != nil {
			L.RaiseError("%s", err)
		}
		table := L.NewTable()
		for i := 0; i < len(res); i++ {
			table.Append(lua.LString(res[i]))
		}
		L.Push(table)
		return 1
	}))
	L.SetGlobal("LGet", L.NewFunction(func(L *lua.LState) int {
		userKey := L.CheckString(1)
		table := L.CheckTable(2)
		indexes := make([]int64, table.Len())
		for i := 0; i < table.Len(); i++ {
			indexes[i] = int64(table.RawGetInt(i + 1).(lua.LNumber))
		}
		luaLogger.Printf("[LGet] userKey: %s, indexes: %+v", userKey, indexes)

		res, err := luaService.list.Get(batch, []byte(userKey), indexes)
		if err != nil {
			L.RaiseError("%s", err)
		}
		table = L.NewTable()
		for index, value := range res {
			table.RawSetInt(int(index), lua.LString(value))
		}
		L.Push(table)
		return 1
	}))
	L.SetGlobal("LContain", L.NewFunction(func(L *lua.LState) int {
		userKey := L.CheckString(1)
		table := L.CheckTable(2)
		values := make([][]byte, table.Len())
		for i := 0; i < table.Len(); i++ {
			values[i] = []byte(table.RawGetInt(i + 1).(lua.LString))
		}
		luaLogger.Printf("[LContain] userKey: %s, values: %+s", userKey, values)

		res, err := luaService.list.Contain(batch, []byte(userKey), values)
		if err != nil {
			L.RaiseError("%s", err)
		}
		table = L.NewTable()
		for i, b := range res {
			table.RawSetInt(i, lua.LBool(b))
		}
		L.Push(table)
		return 1
	}))
	L.SetGlobal("LExist", L.NewFunction(func(L *lua.LState) int {
		userKey := L.CheckString(1)
		luaLogger.Printf("[LExist] userKey: %s", userKey)

		res, err := luaService.list.Exist(batch, []byte(userKey))
		if err != nil {
			L.RaiseError("%s", err)
		}
		L.Push(lua.LBool(res))
		return 1
	}))
	L.SetGlobal("LDelete", L.NewFunction(func(L *lua.LState) int {
		userKey := L.CheckString(1)
		if !found(userKey) {
			L.RaiseError(fmt.Sprintf("can not opt not lock userKey: %s", userKey))
		}
		luaLogger.Printf("[LDelete] userKey: %s", userKey)

		err := luaService.list.Delete(batch, []byte(userKey))
		if err != nil {
			L.RaiseError("%s", err)
		}
		return 0
	}))
	L.SetGlobal("LRem", L.NewFunction(func(L *lua.LState) int {
		userKey := L.CheckString(1)
		if !found(userKey) {
			L.RaiseError(fmt.Sprintf("can not opt not lock userKey: %s", userKey))
		}
		value := L.CheckString(2)
		luaLogger.Printf("[LRem] userKey: %s, value: %s", userKey, value)

		err := luaService.list.Rem(batch, []byte(userKey), []byte(value))
		if err != nil {
			L.RaiseError("%s", err)
		}
		return 0
	}))

	if err := L.DoString(script); err != nil {
		return nil, err
	}

	top := L.GetTop()
	res := make([]string, top)
	for i := 1; i <= top; i++ {
		res[i-1] = L.Get(i).String()
	}

	if err := batch.Commit(&pebble.WriteOptions{Sync: true}); err != nil {
		return nil, err
	}
	return res, nil
}
