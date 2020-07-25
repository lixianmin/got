----

#### 0x01 简介

golang基础库。got英文含【得到】的意思，愿景为避免重复开发



----

#### 0x02 loom

##### 01 高并发字典 map

仿照java.util.concurrent.ConcurrentMap实现的高并发Map类，主要目标为：
1. 代替golang自带的sync.Map，通过sharding提供更高的写并发度
2. 提供像ComputeIfAbsent()这样的延迟初始化方法
3. map的key只支持基本数据类型，包括各种int类型和string



示例：

```go
var m Map

const max = 1000
for i := 0; i < max; i++ {
  m.Put(i, i)
}

const max2 = 2000
for i := max / 2; i < max2; i++ {
  m.ComputeIfAbsent(i, func(key interface{}) interface{} {
    return key.(int) * 2
  })
}
```



----

#### 0x03 sortx

##### 01 二分查找 Search()

目标：在一个有序列表中查找特定的目标值target

算法实现参考《编程珠玑》（人民邮电出版社 第2版）第9.3节 《 大手术 -- 二分搜索》（Page 89）改写



这个算法的特点是：

1. 如果有序列表中存在目标值target，则返回它在有序列表中第1次出现的下标
2. 如果有序列表中不存在目标值target，则返回一个负数下标index，该index的相反数是将target插入到该有序列表中时它应该在的位置



示例：

```go
var list = []int {1, 3, 3, 3, 5, 7, 9, 9, 9, 11}

var target = 9
var index = sortx.Search(len(list), func(i int) bool {
  return list[i] < target
}, func(i int) bool {
  return list[i] == target
})

// 找到的位置，应该是target第一次出现的索引下标
if index != 6 {
  fmt.Println("index should be 6")
}
```



##### 02 Slice逆序 Reverse()

目标：将一个slice逆序

算法实现参考：[golang-reverse-a-arbitrary-slice](https://stackoverflow.com/questions/54858529/golang-reverse-a-arbitrary-slice)



示例：

```go
var list = []int{1, 2, 3, 4, 5, 6}
sortx.Reverse(list)

// 此时list的值为 {6, 5, 4, 3, 2, 1}
```











