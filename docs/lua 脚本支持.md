## SDB 背后的思考 ———— lua 脚本支持
------

lua 脚本是酷炫的，Redis、Nginx 等开源项目都可以内嵌 lua 实现业务逻辑。所以 SDB 也打算支撑 lua 脚本。还好找到了优秀的开源项目：[gopher-lua](https://github.com/yuin/gopher-lua) 。

### 初始化

```go
L := lua.NewState()
defer L.Close()
```

### 方法注册

将 SDB 的方法注册到 lua 脚本中，以 LCount 为例子：

```go
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
```

### 如何保证事务操作？

lua 脚本是灵活的，在上面我们可以写很多的业务逻辑，如果我们在 lua 中用了：LLPush、LRPush 如何能保证多个方法是操作是事务的呢？也比较简单，我们只需要创建一个 batch 对象，执行完我们的 lua 脚本后进行 commit。保证执行每一次脚本是事务的。

### 加锁逻辑？

这其实是最容易被忽略的一点。对比 Redis 来说，由于是单线程的，所以 lua 脚本是不需要考虑锁的。

但是 SDB 不同，SDB 首先是支持多线程的，那么对 userKey 的写操作会进行加锁。 lua 脚本在某种程度上破坏了对单个 userKey 的加锁策略。

首先 lua 脚本得让用户传入会操作的 userKey 列表，然后对每个 userKey 进行加锁。只有获取了所有锁，lua 脚本才能开始运行。

但这是不够的，可能会出现死锁问题。如当一个 lua 脚本操作 a、b 两个 userKey，另一个 lua 脚本操作 b、d 两个 userKey。假设 d 和 a 的锁是同一把。就会出现死锁的问题。

防止死锁的解决方法也很简单，只需要保证先对 lua 脚本的 userKey 按某种顺序依次获取锁既可。伪代码如下：
```go
hashes := make([]bool, len(locker.lockers))
for i := 0; i < len(userKeys); i++ {
	hashes[locker.hash(userKeys[i])] = true
}
for i := 0; i < len(hashes); i++ {
	if hashes[i] {
		locker.lockers[i].Lock()
	}
}
```