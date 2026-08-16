# Go map 原理
## 设计思想
### 根据 key 快速定位 value
### 避免全量遍历
### 均摊查询 $O(1)$

## 核心结构
### hamp
#### count：元素数量，`len(map)`的值
#### B：$bucket 数量 = 2^B$，描述的是当前主 bucket 的数量。要保证是 2 的幂，可以优化成位运算，hash 后低 B 位决定了去哪个 bucket
#### hash0：哈希种子
#### overflow：溢出桶的数量
#### buckets：指向当前桶数组
#### oldbuckets：扩容时指向旧桶数组，翻倍扩容时旧 buckets 数量是新 buckets 的一半，等量扩容时新旧 buckets 数量相同
#### nevacuate：迁移进度，小于它的表示迁移已经完成
### bmap
#### 分布：tophash/key/value连续放，减少内存浪费
#### tophash[8]：hash 后的高 8 位存在这里
#### keys[8]：一 个bucket 里最多 8 个 key，需要进一步过滤
#### value[8]
#### overflow，给哈希冲突兜底，但是溢出桶过多会拖慢访问
#### 先比一字节的 tophash，不同直接到下一个槽位，相同再比较完整的 key

## 定位逻辑
### hash(key)
### 低 B 位选择 bucket
### 高 8 位做 tophash
### tophash 过滤
### 完整 key 比较

## 读取
### nil 或空 map 返回对应类型的零值
### hash 定位 bucket
### 扫描 bucket 和 overflow
### 找不到返回零值，ol=false

## 写入
### nil map 写入 panic
### 先 hash，后标记写状态
### 找到 key 就覆盖
### 找不到 key 就插入空槽位
### 插入前可能触发扩容

## 扩容
### 负载因子 > 6.5：翻倍扩容
### overflow 过多：等量扩容
### hashGrow：切换扩容状态
### growWork：每次帮忙迁移旧桶
### evacuate：搬迁旧桶

## 删除
### delete nil map 是安全的操作
### 定位 bucket 和 slot
### 清理 key/value
### 修改 tophash
### emptyOne：中间空，后面可能还有 key，查找不能停
### emptyRest：后面都是空，查找可以停
### 不会主动缩容

## 遍历
### hiter 迭代器
### 随机 startBucket
### 随机 offset
### 不保证遍历顺序
### range map 不是一致性快照
### 扩容中兼容 oldbuckets / new buckets