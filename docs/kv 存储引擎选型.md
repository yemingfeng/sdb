## SDB 背后的思考 ———— kv 存储引擎选型
------

主流的 kv 存储引擎模型主要分为三大类：

- LSM 树模型
	- 针对写多读少场景，顺序写、读比较慢
- B+ 树模型
	- 针对读多写少场景，随机写、读比较快
- SLM 模型
	- 结合 LSM 树和 B + 树的有点，写入、读取都很快
	- 可参考实现 [lotusdb](https://github.com/flower-corp/lotusdb)

不管哪种存储模型，我们可以将 kv 存储引擎想象成是一个巨大的、持久化的 sortedMap。其提供的接口是：

- get(storeKey)
- set(storeKey, storeValue)
- delete(storeKey)
- scan(minKey, maxKey)

由于 SDB 定位是纯 Go 开发的 NoSQL 数据库，只会兼容纯 Go 开发的 kv 存储库。为了减少工作量，选择集成已有、成熟的 kv 存储引擎，找到了以下可靠的存储引擎：

| 项目 | 介绍 | 存储模型
| ---- | ---- | ---- 
| [pebble](https://github.com/cockroachdb/pebble) | [cockroach](https://github.com/cockroachdb/cockroach) 推出并兼容 rocksdb。 | LSM
| [goleveldb](https://github.com/syndtr/goleveldb) | Go 版 LevelDB | LSM
| [badger](https://github.com/dgraph-io/badger) | Go 版 badger | LSM

目前 SDB 选择了 pebble 作为 kv 存储引擎。PS: 欢迎补充 Go 版本的 kv 存储。