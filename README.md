## [SDB](https://github.com/yemingfeng/sdb) ：纯 Go 开发、数据结构丰富、持久化、简单易用的 NoSQL 数据库
------

### 为什么需要 SDB？

试想以下业务场景：

- 计数服务：对内容的点赞量、播放量进行统计
- 推荐服务：每个用户有一个包含内容和权重推荐列表
- 评论服务：查看内容的评论列表

在传统的做法中，我们会通过 MySQL + Redis 的方式实现。其中 MySQL 提供数据的持久化能力，Redis 提供高性能的读写能力。在这样的架构下会面临以下问题：

- 同时部署 MySQL + Redis，机器成本高
- MySQL + Redis 数据不同步带来的一致性问题
- 随着业务上涨，MySQL 面临海量数据的读写压力

回过头来看上面的需求，我们真正需要的其实是有持久化能力的 Redis。 业内也有解决方案，如：

- [pika](https://github.com/OpenAtomFoundation/pika)
- [kvrocks](https://github.com/apache/incubator-kvrocks)
- [tendis](https://cloud.tencent.com/document/product/1363/50791)
- [memorydb](https://aws.amazon.com/cn/memorydb/)

SDB 也是应对上面问题提供的解决方案。

### RoadMap

- 数据结构
	- [x] List
	- [ ] Linked
	- [ ] String
	- [ ] Set
	- [ ] SortedSet
	- [ ] BitMap
	- [ ] BloomFilter
	- [ ] GeoHash
- [ ] Grpc Server
- [ ] Prometheus 监控
- [ ] kv 存储引擎接入
	- [ ] badger
- [ ] 集群
	- [ ] 主从
	- [ ] 分布式

### sdb 背后的设计

- [kv 存储引擎选型](https://github.com/yemingfeng/sdb/blob/master/docs/kv%20%E5%AD%98%E5%82%A8%E5%BC%95%E6%93%8E%E9%80%89%E5%9E%8B.md)
- [List 数据模型设计](https://github.com/yemingfeng/sdb/blob/master/docs/List%20%E6%95%B0%E6%8D%AE%E6%A8%A1%E5%9E%8B%E8%AE%BE%E8%AE%A1.md)
- [List 锁模型设计](https://github.com/yemingfeng/sdb/blob/master/docs/List%20%E9%94%81%E6%A8%A1%E5%9E%8B%E8%AE%BE%E8%AE%A1.md)
- [lua 脚本支持](https://github.com/yemingfeng/sdb/blob/master/docs/lua%20%E8%84%9A%E6%9C%AC%E6%94%AF%E6%8C%81.md)

### 友链

- [rdb](https://github.com/MoSunDay/rdb)