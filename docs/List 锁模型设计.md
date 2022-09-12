## SDB 背后的思考 ———— List 数据结构锁模型设计
------

[从上一篇中，我们已经设计好了 List 数据结构在 kv 存储中的数据模型。](https://github.com/yemingfeng/sdb/blob/master/docs/List%20%E6%95%B0%E6%8D%AE%E6%A8%A1%E5%9E%8B%E8%AE%BE%E8%AE%A1.md)

这一篇，我们介绍下 List 数据结构中的锁模型设计。

Q：用户同时进行写入操作请求：LPush、LPush、Delete、Rem 等，如何加锁？

A：由于这类请求会操作 metaKey、ttlKey 和 deletedKey。为了保证一致性，SDB 会按照用户写入的 userKey 进行加锁。

SDB 内部维护了多把锁，每个 userKey hash 后会取到对应的锁，然后对该锁进行加锁操作。伪代码如下：

```golang
var lockers []*sync.RWMutex

// 写操作加锁
func wlock(userKey []byte) {
	getLocker(userKey).Lock()
}

// 读操作加锁
func wUnLock(userKey []byte) {
	getLocker(userKey).Unlock()
}

// 根据 userKey 获取锁
func getLocker(userKey []byte) *sync.RWMutex {
	checksum := crc16.Checksum(userKey, crc16.IBMTable)
	return lockers[int(checksum)%len(lockers)]
}
```

Q：为什么不考虑每个 userKey 一把锁，而是多个 userKey 经过 hash 后共用一把锁？

A：如果每个 userKey 一把锁，可能会出现锁太多带来的性能损耗。虽然多个 userKey 经过 hash 后会共用一把锁，但每次用户的写入请求应该是**快速返回**的，写锁锁住的时间应该是不会太长的。

Q：用户的写入操作和 Count、Range 的加锁逻辑是什么？

A：针对 Count，只是读取 meta 信息，不需要做额外的加锁处理。 而 Range 操作是遍历 List 的，为了防止在遍历的时候，对该 List 进行了写入操作带来的数据混乱问题。SDB 对该操作加了读锁，伪代码如下：

```golang
// 写操作加锁
func rlock(userKey []byte) {
	getLocker(userKey).RLock()
}

// 读操作加锁
func rUnLock(userKey []byte) {
	getLocker(userKey).RUnlock()
}
```

Q：deleted 数据回收任务和用户的写入请求加锁逻辑是怎样的？

A：由于 SDB 是采用多版本的设计，用户的写入请求只会操作最新版本的 metaKey 等，而 deleted 数据回收任务不会回收最新版本的数据，所以二者不存在冲突的问题，不需要额外加锁。

Q：ttl 数据回收任务和用户的写入请求加锁逻辑是怎样的？

A：由于 SDB 是采用多版本的设计，用户的写入请求只会操作最新版本的 metaKey 等，而 ttl 数据回收任务可能会回收最新版本的数据，这二者存在以下组合：

- 用户请求 LPush（RPush 同理） + ttl 数据回收任务同时进行
	- 假设 listA 最高版本号是 3 即将过期，meta 信息是：`{count=3, version=3, ttl=1, head=0, tail=3, deleted=false}`，包含了元素是：[a, b, c]。LPush 的元素是：[a, d]。
	- LPush 获取到最高版本号 3，发现未过期（即将过期）。所以会往版本号 3 写入 [a, d] 数据 + meta `{count=5, version=3, ttl=1, head=0, tail=5, deleted=false}` 信息和元素 [a, d]。
	- ttl 数据回收任务只回收了版本号 3 的 metaKey：`{count=3, version=3, ttl=1, head=0, tail=3, deleted=false}` 和数据：[a, b, c]。
	- **总结：这种情况不需要加锁，只需要在每次 LPush 是再保证写入一次 ttlKey 即可。等下一次回收任务还会回收该 listA 的版本号 3 的所有数据。**
- 用户请求 Delete + ttl 数据回收任务同时进行
	- 假设删除 listA 最高版本号是 3，对 listA 进行删除，同时 ttl 数据回收任务发现 listA 即将过期。
	- **总结：这种情况不需要加锁，ttl 数据任务会回收一次 listA 数据，deleted 数据回收也会回收一次 listA 数据。**
	- Rem 同理。