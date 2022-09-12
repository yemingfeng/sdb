package test

import (
	"github.com/yemingfeng/sdb/internal/collection/list"
	"github.com/yemingfeng/sdb/internal/config"
	"github.com/yemingfeng/sdb/internal/service"
	"github.com/yemingfeng/sdb/internal/store"
	"testing"
)

func TestLua(t *testing.T) {
	config := config.NewTestConfig()
	list := list.NewList()
	locker := service.NewLocker(config)
	db := store.NewStore(config)
	luaService := service.NewLuaService(db, locker, list)

	script := `
		res = LExist("h2")
		if (res ~= false)
		then
			LLPush("h2", "v1", "v2", "v3")
			return LCount("h2")
		else
			LRPush("h2", "vv1", "vv2", "vv3")	
			return true 
		end
		return false
	`

	res, err := luaService.Execute([][]byte{
		[]byte("h1"),
		[]byte("h2"),
	}, script)

	t.Logf("res: %+v, err: %+v", res, err)
}
