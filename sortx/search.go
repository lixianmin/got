package sortx

/********************************************************************
created:    2020-07-09
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

// 通用二分查找算法
// 目标：在一个有序列表中查找特定的目标值target
// 算法实现参考《编程珠玑》（人民邮电出版社 第2版）第9.3节 《 大手术 -- 二分搜索》（Page 89）改写
//
// 这个算法的特点是：
// 1. 如果有序列表中存在目标值target，则返回它在有序列表中第1次出现的下标
// 2. 如果有序列表中不存在目标值target，则返回一个负数下标index，该index的相反数是将target插入到该有序列表中时它应该在的位置
//
func Search(count int, less func(i int) bool, equal func(i int) bool) int {
	if count <= 0 {
		return -1
	}

	var i = - 1    // 不变式：data[i]<key
	var j = count  // 不变式：data[j]≥key
	for i+1 != j { // 可以证明：当j≥i+2时[i,j]的范围是在不断缩小的
		var mid = int(uint(i+j) >> 1) // (i+j)/2当数字特别大时有可能溢出
		if less(mid) {
			i = mid
		} else {
			j = mid
		}
	}

	// 查找不成功时如果需要将其插入的话，则无论j==right+1还是data[j]!=key，实际上j实际指向插入的位置
	if j == count || !equal(j) {
		return ^j
	}

	// 查找成功时，j指向key在data[]中第一次出现的位置
	return j
}
