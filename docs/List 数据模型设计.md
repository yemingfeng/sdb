## SDB 背后的思考 ———— List 数据模型设计
------

[我们借助 LSM / B+ 树的实现，已经有了可靠的 kv 存储引擎了](https://github.com/yemingfeng/sdb/blob/master/docs/kv%20%E5%AD%98%E5%82%A8%E5%BC%95%E6%93%8E%E9%80%89%E5%9E%8B.md) 。回顾下我们具备的能力，是一个巨大的、持久化的 sortedMap。提供的接口是：

- get(storeKey)
- set(storeKey, storeValue)
- delete(storeKey)
- scan(minKey, maxKey)

那么，我们如何根据上述接口，实现 List 功能呢？我先抛出 SDB 实现后的结论：

| 接口 | 时间复杂度 | 说明
| ---- | ---- | ---- |
| LPush(userKey, values) | O(1) | 从左边增加元素
| RPush(userKey, values) | O(1) | 从右边增加元素
| LPop(userKey) | O(1) | 从左边弹出元素
| RPop(userKey) | O(1) | 从右边弹出元素
| Count(userKey) | O(1) | 获取元素个数
| List(userKey) | O(n) | 返回某个 listService 所有元素
| Range(userKey, offset, limit) | O(1) | 遍历 List，若 offset < 0，代表反向迭代
| Get(userKey, index) | O(1) | 返回 List 某下标的元素
| Contain(userKey, values) | O(1) | 判断元素是否存在
| Delete(userKey) | O(1) | 删除某个 List
| Rem(userKey, values) | O(n) | 删除 List 中某些元素
| Ttl(userKey, ttl) | O(1) | 设置 ttl
| LInsert(userKey, insertValue, index) | O(n) | 向 List 指定位置左边插入元素
| RInsert(userKey, insertValue, index) | O(n) | 向 List 指定位置右边插入元素

### 存储模型图

<img src="https://raw.githubusercontent.com/yemingfeng/sdb/master/docs/List%20%E6%95%B0%E6%8D%AE%E6%A8%A1%E5%9E%8B.png" />

| storeKey |生成规则 | storeValue | 作用
| ---- | ---- | ---- | ---- |
| metaKey | am:{userKey}:{version} | {head, tail, count, deleted, ttl} | 存储 List 的原始信息，version 在同一个 list 内是递增的
| deletedKey | ad:{userKey}:{version} | - | 若某 List 的 version 被标记为 deleted，则会写入 deletedKey，仅用于快速回收 list 某版本的数据，后面会详细说明
| ttlKey | at:{userKey}:{ttl}:{version} | - | 若某 List 有 ttl，则会写入 ttlKey，仅用于快速回收 List 某版本的数据，后面会详细说明
| seqKey | as:{userKey}:{version}:{seq}:{value} | - | List 每个元素对应一个 seqKey，用于遍历该 List
| valueKey | av:{userKey}:{version}:{value}:{seq} | - | List 每个元素对应一个 valueKey，用于判断某 value 是否在 list 中

### meta 设计

| 字段 | 生成逻辑 | 作用
| ---- | ---- | ---- |
| head | 程序保障 | 指向 List 第一个元素，每向 List 左边增加一个元素会让 head - 1
| tail | 程序保障 | 指向 List 最后一个元素，每向 List 右边增加一个元素会让 tail - 1
| ttl | 用户设置 | 默认值为 0，若 ttl <= 0，则认为该 List 未启用 ttl 功能
| deleted | 用户设置 | 默认值为 false，代表未被删除
| count | 程序保障 | List 的元素个数

假设有一个 listA，其在存储引擎中的数据可能是：

| userKey | version | metaKey | meta 信息
| ---- | ---- | ---- | ---- |
| listA | 0 | `am:listA:0` | `{version=0, count=10, head=0, tail=10, deleted=true, ttl=0}`
| listA | 1 | `am:listA:1` | `{version=1, count=3, head=-2, tail=0, deleted=true, ttl=0}`
| listA | 2 | `am:listA:2` | `{version=2, count=15, head=3, tail=18, deleted=false, ttl=0}`

可以看出，每个 List 的 meta 是包含多版本的。

SDB 只会认为最高版本的 meta 在 ttl 内并且未 deleted 才是有效的。若最高版本的 meta 无效，向列表增加元素时，SDB 会创建一个新的 meta。

非最高版本的数据用户是不能操作的，只会在后面的数据回收任务中回收，后面会详细说明。

### ttl 设计

为了支持 List 过期功能。SDB 在 List meta 中增加了 ttl 字段。 假设 listA 在版本 3 中设置了 ttl=10，则在存储引擎中 listA 的 meta 数据可能是：

| userKey | version | metaKey | meta 信息（省略无用信息）
| ---- | ---- | ---- | ---- |
| listA | 0 | `am:listA:0` | `{ttl=0}`
| listA | 1 | `am:listA:1` | `{ttl=0}`
| listA | 2 | `am:listA:2` | `{ttl=0}`
| listA | 3 | `am:listA:3` | `{ttl=10}`

在用户操作 SDB 时，只会判断 version=3 的 meta 信息，发现无效后，再向 List 加入元素时，SDB 会创建一个新版本的 meta 信息。

假设对 listA 的操作流程如下：

- 向 listA 增加元素 a
- 设置 listA ttl=1（会快速过期）
- 向 listA 增加元素 b
- 设置 listA ttl=1（会快速过期）
- 向 listA 增加元素 c

则 listA 的 meta 信息的数据可能是：

| userKey | version | metaKey | meta 信息（省略无用信息）
| ---- | ---- | ---- | ---- |
| listA | 0 | `am:listA:0` | `{ttl=1}`
| listA | 1 | `am:listA:1` | `{ttl=1}`
| listA | 2 | `am:listA:2` | `{ttl=0}`

其元素的数据可能是：

| userKey | version | 元素 | seqKey | valueKey
| ---- | ---- | ---- | ---- | ---- |
| listA | 0 | a | `as:listA:0:0:a` | `av:listA:0:a:0`
| listA | 1 | b | `as:listA:1:0:b` | `av:listA:1:b:0`
| listA | 2 | c | `as:listA:2:0:c` | `av:listA:2:c:0`

我们可以很快发现 listA 记录中 version=0 和 version=1 相关数据都可以回收。

为了支持快速回收无效数据，SDB 为设置了 ttl 的版本增加了 ttlKey：

| userKey | version | ttlKey
| ---- | ---- | ---- |
| listA | 0 | `at:listA:0`
| listA | 1 | `at:listA:1`

有了 ttlKey 后，SDB 就可以根据它快速扫描出已经 ttl 的 List，用于回收相关数据。

### deleted 设计

同 ttl 设计，不赘述。

### seqKey 和 valueKey 设计

List 中每一个元素都包含 seqKey 和 valueKey，作用已经在上面表述了。 这里补充的是：

- 向 List 左边增加元素时，seq = head - 1
- 向 List 右边增加元素时，seq = tail + 1

假设有一个 listA，它的最高版本是 3，其元素列表如：[a, b, c]，其每个元素在存储引擎中可能是这样的：

| 元素 | seqKey | valueKey
| ---- | ---- | ---- |
| a | `as:listA:3:-1:a` | `av:listA:3:a:-1`
| b | `as:listA:3:0:b` | `av:listA:3:b:0`
| c | `as:listA:3:1:c` | `av:listA:3:c:1`

### 如何在指定位置插入元素？

假设有一个 listA，元素为[a, b, c]，其存储数据为：

| 元素 | seqKey | valueKey
| ---- | ---- | ---- |
| a | `as:listA:0:-1:a` | `av:listA:0:a:-1`
| b | `as:listA:0:0:b` | `av:listA:0:b:0`
| c | `as:listA:0:1:c` | `av:listA:0:c:1`

现在要在 b 后面增加一个元素 d，则存储数据变为：

| 元素 | seqKey | valueKey
| ---- | ---- | ---- |
| a | `as:listA:0:-1:a` | `av:listA:0:a:-1`
| b | `as:listA:0:0:b` | `av:listA:0:b:0`
| d | `as:listA:0:1:d` | `av:listA:0:d:1`
| c | `as:listA:0:2:c` | `av:listA:0:c:2`

意味着所有的元素都会挪动，时间复杂度较高。Rem 接口同理。