package osx

import "os"

/********************************************************************
created:    2020-01-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

// IsPathExist 判断所给路径文件/文件夹是否存在
func IsPathExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}

	return true
}

// 使用 os.MkdirAll()代替
// perm 常为 0600 -> r:4, w:2, x:1 , 0600的意思是 owner可rw，参考：https://chmodcommand.com/chmod-600/
//func EnsureDir(path string, perm os.FileMode) error {
//	if !IsPathExist(path) {
//		return os.MkdirAll(path, perm)
//	}
//
//	return nil
//}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}
