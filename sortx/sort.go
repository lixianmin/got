package sortx

/********************************************************************
created:    2020-07-09
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

// 二分查找
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
