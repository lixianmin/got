package loom

//
//import (
//	"fmt"
//	"sync"
//	"sync/atomic"
//	"unsafe"
//)
//
///********************************************************************
//created:    2020-07-13
//author:     lixianmin
//
//deprecated: 2020-07-13
//	1. 这个类不再维护，因为sync.Map已经足够了
//
//////////////////////////////////////////////////////////////////////
//仿jdk1.8之前ConcurrentHashMap中segment实现的Map类，主要目标为：
//1. 提供更高的写并发度
//2. 提供像ComputeIfAbsent()这样的延迟初始化方法
//
// 问题: 是否将loom.Map分uint32, uint64与string，分类型初始化？
//	1. 这对速度和内存占用可能会是一个bonus，我们通常只需要value是any，而key不需要太复杂的类型（本来我们也没支持几种类型）；
//	2. 这会导致同时只能存储单一的key类型，比如只能存储int或string。不过，大部分情况下我们的key的确是单一类型；
//	3. 编译器不能帮助发现错误，这有可能会导致运行时的bug
//	4. 如果分为IntMap与StringMap的话，则会导致相同的代码写好两遍，并给使用和维护带来不便
//	5. 这个问题也许可以靠10年以后的泛型解决
//
// 问题: 减少默认sharding数？同时增加一个只允许调用一次SetSharding(count)方法
//	1. 跟直接加一个NewMap(shardingCount)初始化方法相比孰优孰劣？减小初始化大小会降低内存使用量（为什么印象里实测并没有减少内存使用量？）
//	2. 会增加使用复杂度，因为大多数人都只使用默认的设置；
//	3. 强迫使用NewMap(shardingCount)方法增加使用门槛；
//
//Copyright (C) - All Rights Reserved
//*********************************************************************/
//
//// 用于记录快照数据的共享池
//var snapshotPool = newPool(8)
//var mapSharding = NewSharding()
//
//type ShardTable map[any]any
//
//type shardItem struct {
//	sync.RWMutex
//	items ShardTable
//}
//
//type Map struct {
//	lock sync.Mutex
//	data *[]*shardItem
//	size int64
//}
//
//// Put 如果已经存在了相同key的value，则覆盖找返回以前存在的那一个值；否则返回nil
//func (my *Map) Put(key any, value any) any {
//	var shard, normalizedKey = my.getShard(key)
//	var last any
//	var has = false
//	shard.Lock()
//	{
//		last, has = shard.items[normalizedKey]
//		if has {
//			shard.items[normalizedKey] = value
//		} else {
//			shard.items[normalizedKey] = value
//			my.addSize(1)
//		}
//	}
//	shard.Unlock()
//	return last
//}
//
//// Remove 如果存在，则删除，并返回该值
//func (my *Map) Remove(key any) any {
//	var shard, normalizedKey = my.getShard(key)
//	var last any
//	var has = false
//	shard.Lock()
//	{
//		last, has = shard.items[normalizedKey]
//		if has {
//			delete(shard.items, normalizedKey)
//			my.addSize(-1)
//		}
//	}
//	shard.Unlock()
//	return last
//}
//
//// Get1 如果map中存在，则返回；否则返回nil
//func (my *Map) Get1(key any) any {
//	var shard, normalizedKey = my.getShard(key)
//	var last, _ = my.getInner(shard, normalizedKey)
//	return last
//}
//
//func (my *Map) Get2(key any) (any, bool) {
//	var shard, normalizedKey = my.getShard(key)
//	return my.getInner(shard, normalizedKey)
//}
//
//func (my *Map) getInner(shard *shardItem, key any) (any, bool) {
//	var last any
//	var has = false
//	shard.RLock()
//	{
//		last, has = shard.items[key]
//	}
//	shard.RUnlock()
//	return last, has
//}
//
//// PutIfAbsent 这其实是一种get命令：如果key对应的value已经存在，则返回存在的value，不进行替换；如果不存在，就添加key和value，然后返回nil
//func (my *Map) PutIfAbsent(key any, value any) any {
//	var shard, normalizedKey = my.getShard(key)
//	var last, has = my.getInner(shard, normalizedKey)
//	if has {
//		return last
//	}
//
//	shard.Lock()
//	{
//		last, has = shard.items[normalizedKey]
//		if !has {
//			shard.items[normalizedKey] = value
//			my.addSize(1)
//		}
//	}
//	shard.Unlock()
//	return last
//}
//
//// ComputeIfAbsent 如果原来存在，则返回原来的值；否则使用creator创建一个新值，放到到map中，则返回它
//func (my *Map) ComputeIfAbsent(key any, creator func(key any) any) any {
//	var shard, normalizedKey = my.getShard(key)
//	var last, has = my.getInner(shard, normalizedKey)
//	if has {
//		return last
//	}
//
//	shard.Lock()
//	defer shard.Unlock() // 用defer是因为不知道creator会不会panic
//
//	// 加x锁后需要重新测试有没有数据
//	// 如果creator=nil，也就不会重新生成了
//	last, has = shard.items[normalizedKey]
//	if has || creator == nil {
//		return last
//	}
//
//	// 如果没有，则创建一个放入到容器中
//	var item = creator(key)
//	shard.items[normalizedKey] = item
//	my.addSize(1)
//	return item
//}
//
//// 感觉这个方法不应该被开出来，不记得当时写这个方法用来处理哪个项目的问题的，先关闭，有问题再说
////// 为什么会有这么奇怪的一个方法？有时，我们需要在锁定某个key的情况下执行某些操作，防止在操作的过程中该key被插入导致不一致性
////func (my *Map) WithLock(key any, handler func(table ShardTable)) {
////	var shard, _ = my.getShard(key)
////	shard.Lock()
////	defer shard.Unlock() // 用defer是因为不知道handler会不会panic
////	handler(shard.items)
////}
//
//// Range 使用数据快照支持遍历过程，因为无法保证遍历过程中回调方法不调用Add, Remove等方法，必须得避免死锁
//func (my *Map) Range(handler func(key any, value any)) {
//	if handler == nil {
//		return
//	}
//
//	var pData = my.getShardData()
//	if pData != nil {
//		var list = snapshotPool.Get()
//		var data = *pData
//		for i := range data {
//			var shard = data[i]
//			// 获取快照
//			shard.RLock()
//			for k, v := range shard.items {
//				list = append(list, k, v)
//			}
//			shard.RUnlock()
//
//			// 调用回调方法
//			var size = len(list)
//			for i := 0; i < size; i += 2 {
//				var k, v = list[i], list[i+1]
//				safeRangeHandler(handler, k, v)
//			}
//			// 重置list
//			list = list[:0]
//		}
//
//		snapshotPool.Put(list)
//	}
//}
//
//func safeRangeHandler(f func(key any, value any), key any, value any) {
//	defer func() {
//		if r := recover(); r != nil {
//			fmt.Println(r)
//		}
//	}()
//
//	f(key, value)
//}
//
//func (my *Map) Size() int {
//	return int(atomic.LoadInt64(&my.size))
//}
//
//// 调整my.size时使用atomic.AddInt64()而不是直接my.size += delta，参考sync.Once的doSlow()中对done字段的修改
//func (my *Map) addSize(delta int64) {
//	atomic.AddInt64(&my.size, delta)
//}
//
//func (my *Map) getShard(key any) (*shardItem, any) {
//	var pData = my.getShardData()
//	var index, normalizedKey = mapSharding.GetShardingIndex(key)
//	var shard = (*pData)[index]
//	return shard, normalizedKey
//}
//
//func (my *Map) getShardData() *[]*shardItem {
//	var pData = (*[]*shardItem)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&my.data))))
//	if pData == nil {
//		pData = my.getDataSlow()
//	}
//
//	return pData
//}
//
//// 将slow方法提取出来，减小主方法体的大小，提高主方法体inline的可能性
//func (my *Map) getDataSlow() *[]*shardItem {
//	my.lock.Lock()
//	var pData = my.data
//	if pData == nil {
//		var shardingCount = mapSharding.GetShardingCount()
//		var slice = make([]*shardItem, shardingCount)
//		for i := 0; i < shardingCount; i++ {
//			var item = &shardItem{items: make(ShardTable, 4)}
//			slice[i] = item
//		}
//
//		pData = &slice
//		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&my.data)), unsafe.Pointer(pData))
//	}
//	my.lock.Unlock()
//	return pData
//}
