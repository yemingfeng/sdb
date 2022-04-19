## [SDB](https://github.com/yemingfeng/sdb) ：纯 golang 开发、分布式、数据结构丰富、持久化、简单易用的 NoSQL 数据库
------

### 为什么需要 SDB？

试想以下业务场景：

- 计数服务：对内容的点赞、播放等数据进行统计
- 评论服务：发布评论后，查看某个内容的评论列表
- 推荐服务：每个用户有一个包含内容和权重推荐列表

以上几个业务场景，都可以通过 MySQL + Redis 的方式实现。

MySQL 在这个场景中充当了持久化的能力，Redis 提供了在线服务的读写能力。

上面的存储要求是：持久化 + 高性能 + 数据结构丰富。能不能只使用一个存储就满足上面的场景呢？

答案是：非常少的。有些数据库要么是支持的数据结构不够丰富，要么数据结构不够丰富，要么是接入成本太高。。。。。。

**为了解决上述问题，SDB 产生了。**

------

### SDB 简单介绍

- 纯 golang 开发，核心代码不超过 1k，代码易读
- 数据结构丰富
    - [string](https://github.com/yemingfeng/sdb/blob/master/examples/string.go)
    - [list](https://github.com/yemingfeng/sdb/blob/master/examples/list.go)
    - [set](https://github.com/yemingfeng/sdb/blob/master/examples/set.go)
    - [sorted set](https://github.com/yemingfeng/sdb/blob/master/examples/sorted_set.go)
    - [bloom filter](https://github.com/yemingfeng/sdb/blob/master/examples/bloom_filter.go)
    - [hyper log log](https://github.com/yemingfeng/sdb/blob/master/examples/hyper_log_log.go)
    - [bitset](https://github.com/yemingfeng/sdb/blob/master/examples/bitset.go)
    - [map](https://github.com/yemingfeng/sdb/blob/master/examples/map.go)
    - [geo hash](https://github.com/yemingfeng/sdb/blob/master/examples/geo_hash.go)
    - [pub sub](https://github.com/yemingfeng/sdb/blob/master/examples/pub_sub.go)
- 持久化
    - 兼容 [pebble](https://github.com/cockroachdb/pebble)
      、[leveldb](https://github.com/syndtr/goleveldb)
      、[badger](https://github.com/dgraph-io/badger) 存储引擎
- 监控
    - 支持 prometheus + grafana 监控方案
- cli
    - 简单易用的 [cli](https://github.com/yemingfeng/sdb-cli)
- 分布式
    - 基于 raft 实现了主从架构

------

### 架构

<img alt="architecture" src="https://github.com/yemingfeng/sdb/raw/master/docs/architecture.png" width="50%" height="50%"/>

- 基于 raft 协议实现了主从架构。
- 当发起写入请求时，可以连接任意节点，若该节点如果是主节点，则直接处理请求，否则该请求将转发到主节点。
- 当发起读取请求时，可以连接任意节点，该节点会直接处理请求。

------

### 快速使用

#### 编译 protobuf

由于使用了 protobuf，该项目并没有将 protobuf 生成的 go 文件上传到 github。 需要手动触发编译 protobuf 文件

```shell
sh ./scripts/build_protobuf.sh
```

#### 启动主节点

```shell
sh ./scripts/start_sdb.sh
```

#### 启动从节点 1
```shell
sh ./scripts/start_slave1.sh
```

#### 启动从节点 2
```shell
sh ./scripts/start_slave2.sh
```

**默认使用 pebble 存储引擎。**

#### 使用 [cli](https://github.com/yemingfeng/sdb-cli)

#### 客户端使用

```go
package main

import (
	"github.com/yemingfeng/sdb/internal/util"
	pb "github.com/yemingfeng/sdb/pkg/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var clientLogger = util.GetLogger("client")

func main() {
	conn, err := grpc.Dial(":10000", grpc.WithInsecure())
	if err != nil {
		clientLogger.Printf("faild to connect: %+v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	
	c := pb.NewSDBClient(conn)
	setResponse, err := c.Set(context.Background(),
		&pb.SetRequest{Key: []byte("hello"), Value: []byte("world")})
	clientLogger.Printf("setResponse: %+v, err: %+v", setResponse, err)
	getResponse, err := c.Get(context.Background(),
		&pb.GetRequest{Key: []byte("hello")})
	clientLogger.Printf("getResponse: %+v, err: %+v", getResponse, err)
}
```

------

### 性能测试

测试脚本：[benchmark](https://github.com/yemingfeng/sdb/blob/master/examples/benchmark_sdb.go)

测试机器：MacBook Pro (13-inch, 2016, Four Thunderbolt 3 Ports)

处理器：2.9GHz 双核 Core i5

内存：8GB

**测试结果： peek QPS > 10k，avg QPS > 8k，set avg time < 80ms，get avg time <
0.2ms**

<img alt="benchmark" src="https://github.com/yemingfeng/sdb/raw/master/docs/benchmark.png" width="50%" height="50%" />

------

### 规划

- [x] 编写接口文档
- [x] 实现更多的 api (2021.12.30)
    - [x] String
        - [x] SetNX
        - [x] GetSet
        - [x] MGet
        - [x] MSet
    - [x] List
        - [x] LMembers
        - [x] LLPush
    - [x] Set
        - [x] SMembers
    - [x] Sorted Set
        - [x] ZMembers
- [x] 支持更丰富的数据结构 (2021.01.20)
    - [x] bitset
    - [x] map
    - [x] geo hash
- [x] [sdb-cli](https://github.com/yemingfeng/sdb-cli) (2021.03.10)
- [x] 主从架构
- [ ] 编写 sdb kv 存储引擎

------

### 接口文档

#### string

接口 | 参数 | 描述
---- | --- | ---
Set | key, value | 设置 kv
MSet | keys, values | 设置一组 kv
SetNX | key, value | 当 key 不存在时，设置 value
Get | key | 获取 key 对应的 value
MGet | keys | 获取一组 key 对应的 value
Del | key | 删除一个 key
Incr | key, delta | 对 key 进行加 delta 操作，如果 value 不为数字，则抛出异常。如果 value 不存在，则 value = delta

#### list

接口 | 参数 | 描述
---- | --- | ---
LRPush | key, values | 从 key 数组后面追加 values
LLPush | key, values | 从 key 数组前面追加 values
LPop | keys, values | 删除 key 数组中所有的 values 元素
LRange | key, offset, limit | 按数组顺序遍历 key，从 0 开始。如果 offset = -1，则从后向前遍历
LExist | key, values | 判断 values 是否存在 key 数组中
LDel | key | 删除某个 key 数组
LCount | key | 返回 key 数组中的元素个数，时间复杂度较高，**不推荐使用**
LMembers | key | 按数组顺序遍历 key。时间复杂度较高，**不推荐使用**

#### set

接口 | 参数 | 描述
---- | --- | ---
SPush | key, values | 把 values 加到 key 集合中
SPop | keys, values | 删除 key 集合中所有的 values 元素
SExist | key, values | 判断 values 是否存在 key 集合中
SDel | key | 删除某个 key 集合
SCount | key | 返回 key 集合中的元素个数，时间复杂度较高，**不推荐使用**
SMembers | key | 按 value 大小遍历 key。时间复杂度较高，**不推荐使用**

#### sorted set

接口 | 参数 | 描述
---- | --- | ---
ZPush | key, tuples | 把 values 加到 key 有序集合中，按 tuple.score 从小到大排序
ZPop | keys, values | 删除 key 有序集合中所有的 values 元素
ZRange | key, offset, limit | 按 score 大小，从小到大遍历 key。如果 offset = -1，则按 score 从大到小开始遍历
ZExist | key, values | 判断 values 是否存在 key 有序集合中
ZDel | key | 删除某个 key 有序集合
ZCount | key | 返回 key 有序集合中的元素个数，时间复杂度较高，**不推荐使用**
ZMembers | key | 按 score 大小，从小到大遍历 key。时间复杂度较高，**不推荐使用**

#### bloom filter

接口 | 参数 | 描述
---- | --- | ---
BFCreate | key, n, p | 创建 bloom filter，n 是元素个数，p 是误判率
BFDel | key | 删除某个 key bloom filter
BFAdd | key, values | 把 values 加入到 bloom filter 中。当 bloom filter 未创建时，将抛出异常
BFExist | key, values | 判断 values 是否存在 key bloom filter 中

#### hyper log log

接口 | 参数 | 描述
---- | --- | ---
HLLCreate | key | 创建 hyper log log
HLLDel | key | 删除某个 key hyper log log
HLLAdd | key, values | 把 values 加入到 hyper log log 中。当 hyper log log 未创建时，将抛出异常
HLLCount | key | 获取某个 hyper log log 的去重元素个数

#### bitset

接口 | 参数 | 描述
---- | --- | ---
BSDel | key | 删除某个 key bitset
BSSetRange | key, start, end, value | 将 key [start, end) 范围的 bit 设置为 value
BSMSet | key, bits, value | 将 key bits 设置为 value
BSGetRange | key, start, end | 获取 key [start, end) 范围的 bit
BSMGet | key, bits | 获取 key bits 的 bit
BSMCount | key | 获取 key bit = 1 的个数
BSMCountRange | key, start, end | 获取 key [start, end) bit = 1 的个数

#### map

接口 | 参数 | 描述
---- | --- | ---
MPush | key, pairs | 把 pairs KV 对加到 key map 中
MPop | key, keys | 删除 key map 中所有的 keys 元素
MExist | key, keys | 判断 keys 是否存在 key map 中
MDel | key | 删除某个 key map
MCount | key | 返回 key map 中的元素个数，时间复杂度较高，**不推荐使用**
MMembers | key | 按 pair.key 大小遍历 pair。时间复杂度较高，**不推荐使用**

#### geo hash

接口 | 参数 | 描述
---- | --- | ---
GHCreate | key, precision | 创建 geo hash，precision 代表精度。
GHDel | key | 删除某个 geo hash
GHAdd | key, points | 将 points 加入到 geo hash 中，point 中的 id 作为唯一标识
GHPop | key, ids | 删除某 points
GHGetBoxes | key, point | 返回和某 point 在 key geo hash 相同 box 的 point 列表，会按照距离从小到大排序
GHGetNeighbors | key, point | 返回在 key geo hash 中距离 point 最近的 point 列表，会按照距离从小到大排序
GHCount | key | 返回 key geo hash 中的元素个数，时间复杂度较高，**不推荐使用**
GHMembers | key | 返回 key geo hash 中所有的 point 列表。时间复杂度较高，**不推荐使用**

#### page

接口 | 参数 | 描述
---- | --- | ---
PList | dataType, key, offset, limit | 查询某个 [dataType](https://github.com/yemingfeng/sdb-protobuf/blob/master/data_type.proto) 下已有的元素。key 如果不为空，则获取该 dataType 下 key 的元素。

#### pub sub

接口 | 参数 | 描述
---- | --- | ---
Subscribe | topic | 订阅某个 topic
Publish | topic, payload | 向某个 topic 发布 payload

#### cluster

接口 | 参数 | 描述
---- | --- | ---
CInfo |  | 获取集群中的节点信息

------

### 监控

#### 安装 docker 版本 grafana、prometheus（可跳过）

- 启动 [scripts/run_monitor.sh](https://github.com/yemingfeng/sdb/blob/master/scripts/run_monitor.sh)

#### 配置 grafana

- 打开 grafana：http://localhost:3000 （注意替换 ip 地址）
- 新建 prometheus datasources：http://host.docker.internal:9090 （如果使用 docker 安装则为这个地址。如果
  host.docker.internal
  无法访问，就直接替换 [prometheus.yml](https://github.com/yemingfeng/sdb/blob/master/scripts/prometheus.yml)
  文件的 host.docker.internal 为自己的 ip 地址就行）
- 将 [scripts/dashboard.json](https://github.com/yemingfeng/sdb/blob/master/scripts/dashboard.json)
  文件导入 grafana dashboard

最终效果可参考：性能测试的 grafana 图

------

### [配置参数](https://github.com/yemingfeng/sdb/blob/master/configs/config.yml)

参数名 | 含义 | 默认值
---- | --- | ---
store.engine | 存储引擎，可选 pebble、level、badger | pebble
store.path | 存储目录 | ./master/db/
server.grpc_port | grpc 监听的端口 | 10000
server.http_port | http 监控的端口，供 prometheus 使用 | 11000
server.rate | 每秒 qps 的限制 | 30000
cluster.path | raft 日志存储的目录 | ./master/raft
cluster.node_id | raft 协议标识的 node id，得唯一 | 1
cluster.address | raft 通讯的地址 | 127.0.0.1:12000
cluster.master | 现有集群中的主节点地址，通过主节点暴露的 http【http_port】 接口进行加入，若是新集群，则为空 | 
cluster.timeout | raft 协议 apply timeout，单位是 ms | 10000
cluster.join | 作为从节点，是否要加入到主节点；首次加入需要置为 true，加入后再次启动需置为 false | false

------

### SDB 原理之——存储引擎选型

SDB 项目最核心的问题是数据存储方案的问题。

首先，我们不可能手写一个存储引擎。这个工作量太大，而且不可靠。 我们得在开源项目中找到适合 SDB 定位的存储方案。

SDB 需要能够提供高性能读写能力的存储引擎。 单机存储引擎方案常用的有：B+ 树、LSM 树、B 树等。

还有一个前置背景，golang 在云原生的表现非常不错，而且性能堪比 C 语言，开发效率也高，所以 SDB 首选使用纯 golang 进行开发。

那么现在的问题变成了：找到一款纯 golang 版本开发的存储引擎，这是比较有难度的。收集了一系列资料后，找到了以下开源方案：

- LSM 树
    - [go-leveldb](https://github.com/golang/leveldb/) ：是一个 unstable 的项目，无法使用
    - [syndtr-goleveldb](https://github.com/syndtr/goleveldb)
    - [badger](https://github.com/dgraph-io/badger)
    - [pebble](https://github.com/cockroachdb/pebble)
- B+ 树
    - [boltdb-bolt](https://github.com/boltdb/bolt) ：是废弃的项目，无法使用
    - [etcd-bolt](https://github.com/etcd-io/bbolt) ：主要是用于分布式环境下的数据同步，无法应对高并发的数据读写

综合来看，golangdb、badger、pebble 这三款存储引擎都是很不错的。

为了兼容这三款存储引擎，SDB 提供了抽象的[接口](https://github.com/yemingfeng/sdb/blob/master/internal/store/store.go)
，进而适配这三个存储引擎。

### SDB 原理之——数据结构设计

SDB 已经通过之前的三款存储引擎解决了数据存储的问题了。 但如何在 KV 的存储引擎上支持丰富的数据结构呢？

以 pebble 为例子，首先 pebble 提供了以下的接口能力：

- set(k, v)
- get(k)
- del(k)
- batch
- iterator

接下来，我以支持 List 数据结构为例子，剖析下 SDB 是如何通过 pebble 存储引擎支持 List 的。

List 数据结构提供了以下接口：LRPush、LLPush、LPop、LExist、LRange、LCount。

如果一个 List 的 key 为：[hello]，该 List 的列表元素有：[aaa, ccc, bbb]，那么该 List 的每个元素在 pebble 的存储为：

pebble key | pebble value
---- | ---
l/hello/{unique_ordering_key1} | aaa
l/hello/{unique_ordering_key2} | ccc
l/hello/{unique_ordering_key3} | bbb

List 元素的 pebble key 生成策略：

- 数据结构前缀：List 都以 **l** 字符为前缀，Set 是以 **s** 为前缀...
- List key 部分：List 的 key 为 hello
- unique_ordering_key：生成方式是通过雪花算法实现的，雪花算法保证局部自增
- pebble value 部分：List 元素真正的内容，如 aaa、ccc、bbb

为什么这么就能保证 List 的插入顺序呢？

这是因为 pebble 是 LSM 的实现，内部使用 key 的字典序排序。为了保证插入顺序，SDB 在 pebble key 中增加了 unique_ordering_key
作为排序的依据，从而保证了插入顺序。

有了 pebble key 的生成策略，一切都变得简单起来了。我们看看 LRPush、LLPush、LPop、LRange 的核心逻辑：

#### LRPush

```go
func LRPush(key []byte, values [][]byte) (bool, error) {
	batchAction := store.NewBatchAction()
	defer batchAction.Close()

	for _, value := range values {
		batchAction.Set(generateListKey(key, util.GetOrderingKey()), value)
	}

	return batchAction.Commit()
}
```

#### LLPush

LLPush 的逻辑和 LRPush 的逻辑非常类似，不同的地方在于，只要将 {unique_ordering_key} 取负数，变成最小值就可以了。 为了保证 values 内部有序，所以还得 -
index。 逻辑如下：

```go
func LLPush(key []byte, values [][]byte) (bool, error) {
	batch := store.NewBatch()
	defer batch.Close()

	for i, value := range values {
		batch.Set(generateListKey(key, -(math.MaxInt64 - util.GetOrderingKey())), value)
	}

	return batch.Commit()
}
```

#### LPop

在写入到 pebble 的时候，key 的生成是通过 unique_ordering_key 的方案。 无法直接在 pebble 中找到 List 的元素在 pebble
key。在删除一个元素的时候，需要遍历 List 的所有元素，找到 value = 待删除的元素，然后进行删除。核心逻辑如下：

```go
func LPop(key []byte, values [][]byte) (bool, error) {
	batchAction := store.NewBatchAction()
	defer batchAction.Close()

	store.Iterate(&store.IteratorOption{Prefix: generateListPrefixKey(key)},
		func(key []byte, value []byte) {
			for i := range values {
				if bytes.Equal(values[i], value) {
					batchAction.Del(key)
				}
			}
		})

	return batchAction.Commit()
}
```

#### LRange

和删除逻辑类似，通过 iterator
接口进行遍历。 [这里对反向迭代做了额外的支持](https://github.com/yemingfeng/sdb/blob/master/internal/store/store.go#L25)
允许 offset 传入 -1，代表从后进行迭代。

```go
func LRange(key []byte, offset int32, limit uint32) ([][]byte, error) {
	index := int32(0)
	res := make([][]byte, limit)
	store.Iterate(&engine.PrefixIteratorOption{
		Prefix: generateListPrefixKey(key), Offset: offset, Limit: limit},
		func(key []byte, value []byte) {
			res[index] = value
			index++
		})
	return res[0:index], nil
}
```

以上就实现了对 List 的数据结构的支持。

其他的数据结构大体逻辑类似，其中 [sorted_set](https://github.com/yemingfeng/sdb/blob/master/internal/service/sorted_set.go)
更加复杂些。可以自行查看。

#### LPop 优化

聪明的大家可以看出，LPop 的逻辑在数据量很大的情况下，非常耗性能。是因为我们在存储引擎中是无法知道 value 对应的 key 的，需要需要将 List 中的元素全部 load
出来后，挨个判断，才能进行删除。

为了降低时间复杂度，提高性能。 还是以 List: [hello] -> [aaa, ccc, bbb] 为例子。存储模型将改成如下：

正排索引结构【不变】：

pebble key | pebble value
---- | ---
l/hello/{unique_ordering_key1} | aaa
l/hello/{unique_ordering_key2} | ccc
l/hello/{unique_ordering_key3} | bbb

辅助索引结构

pebble key | pebble value
---- | ---
l/hello/aaa/{unique_ordering_key1} | aaa
l/hello/ccc/{unique_ordering_key2} | ccc
l/hello/bbb/{unique_ordering_key3} | bbb

有了这个辅助索引后，我们可以通过前缀检索的方式，判断 List 是否存在某个 value 的元素。从而降低时间复杂度，提高性能。 这里面还需要在写入元素时，将辅助索引写入，所以核心逻辑将改成：

```go
func LRPush(key []byte, values [][]byte) (bool, error) {
	batch := store.NewBatch()
	defer batch.Close()

	for _, value := range values {
		id := util.GetOrderingKey()
		batch.Set(generateListKey(key, id), value)
		batch.Set(generateListIdKey(key, value, id), value)
	}

	return batch.Commit()
}

func LLPush(key []byte, values [][]byte) (bool, error) {
	batch := store.NewBatch()
	defer batch.Close()

	for i, value := range values {
		id := -(math.MaxInt64 - util.GetOrderingKey())
		batch.Set(generateListKey(key, id), value)
		batch.Set(generateListIdKey(key, value, id), value)
	}

	return batch.Commit()
}

func LPop(key []byte, values [][]byte) (bool, error) {
	batch := store.NewBatch()
	defer batch.Close()

	for i := range values {
		store.Iterate(&engine.PrefixIteratorOption{Prefix: generateListIdPrefixKey(key, values[i])},
			func(storeKey []byte, storeValue []byte) {
				if bytes.Equal(storeValue, values[i]) {
					batch.Del(storeKey)

					infos := strings.Split(string(storeKey), "/")
					id, _ := strconv.ParseInt(infos[len(infos)-1], 10, 64)
					batch.Del(generateListKey(key, id))
				}
			})
	}

	return batch.Commit()
}
```

------

### SDB 原理之——关系模型到 KV 模型的映射

有了上面的方案，大概知道了 KV 型存储引擎如何支持数据结构。但这种方式很粗暴，无法通用化。

参考了 [TiDB 的设计](https://pingcap.com/zh/blog/tidb-internal-2) ，SDB 做了一层关系模型到 KV 结构的设计。

在 SBD 中，数据由 Collection 和 Row 构造。 其中：

- [Collection](https://github.com/yemingfeng/sdb/blob/master/internal/store/collection.go#L30)
  类似数据库的一张表，是逻辑概念。每个 dataType(如 List) 对应一个 Collection。一个 Collection 包含多个 Row。
- 一个 Row 包含唯一键：key、id、value、indexes，**是真正存储于 KV 存储的数据**。每行 row 以 rowKey 作为唯一值，rowKey
  = `{dataType} + {key} + {id}`
- 每个 row 包含 N 个索引，每个索引以 indexKey 作为唯一值，indexKey
  = `{dataType} + {key} + idx_{indexName} + {indexValue} + {id}`

以 ListCollection 为例子，该 List 的 key 为 [l1]，假设该 Collection 有 4 行 Row，每行 Row 都有 value 和 score 的索引

那么每行 Row 结构如下：

```yaml
 { {key: l1}, {id: 1.1}, {value: aaa}, {score: 1.1}, indexes: [ {name: "value", value: aaa}, {name: "score", value: 1.1} ] }
 { {key: l1}, {id: 2.2}, {value: bbb}, {score: 2.2}, indexes: [ {name: "value", value: bbb}, {name: "score", value: 2.2} ] }
 { {key: l1}, {id: 3.3}, {value: ccc}, {score: 3.3}, indexes: [ {name: "value", value: ccc}, {name: "score", value: 3.3} ] } 
 { {key: l1}, {id: 4.4}, {value: aaa}, {score: 4.4}, indexes: [ {name: "value", value: aaa}, {name: "score", value: 4.4} ] }
```

以 id = 1.1 的 Row 为例子，dataType = 1，rowKey = `1/l1/1.1`，valueIndexKey =
`1/l1/idx_value/aaa/1.1`, scoreIndexKey = `1/l1/idx_score/1.1/1.1` 写入的数据为：

```yaml
    rowKey: 1/l1/1.1 -> { {key: l1}, {id: 1.1}, {value: aaa}, {score: 1.1}, indexes: [ {name: "value", value: aaa}, {name: "score", value: 1.1} ] }
    valueIndexKey: 1/l1/idx_value/aaa/1.1, -> 1/l1/1.1
    scoreIndexKey: 1/l1/idx_score/1.1/1.1 -> 1/l1/1.1
```

如此，便将数据结构、KV 存储、关系模型打通。

------

### SDB 原理之——通讯协议方案

解决完了存储和数据结构的问题后，SDB 面临了【最后一公里】的问题是通讯协议的选择。

SDB 的定位是支持多语言的，所以需要选择支持多语言的通讯框架。

grpc 是一个非常不错的选择，只需要使用 SDB proto 文件，就能通过 protoc 命令行工具自动生成各种语言的客户端，解决了需要开发不同客户端的问题。

------

### SDB 原理之——主从架构

解决完上述的所有问题，SDB 已经是一个可靠的单机 NoSQL 数据库了。

为了更近一步优化 SDB，为 SDB 加入了主从架构。

golang 语言下的 raft 协议有三种选择，分别是：
- [hashicorp/raft](https://github.com/hashicorp/raft)
- [etcd/raft](https://github.com/etcd-io/etcd/blob/main/raft/raft.go)
- [dragonboat/raft](https://github.com/lni/dragonboat)

为了支持国人项目，选择了 dragonboat。PS：感谢国人！

使用了 raft 协议后，整个集群的吞吐量下降了一点点，但有了水平拓展的读能力，收益还是很明显的。

------

### 版本更新记录

#### v1.7.0

- [commit](https://github.com/yemingfeng/sdb/commit/5e29f5bf50847898cbffa9046df75f2f4fa3ffb6)
  使用分片的方式存储 bitset，bitset 不再需要初始化，有了【自动扩容】的功能

### **感谢开源的力量，这里就不一一列举了，请大家移步 [go.mod](https://github.com/yemingfeng/sdb/blob/master/go.mod)**
