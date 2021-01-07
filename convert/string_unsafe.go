// +build !appengine

package convert

import "unsafe"

/********************************************************************
created:    2021-01-06
author:     lixianmin

this file is extracted from go-redis/v8/internal/unsafe.go
*********************************************************************/

// String converts byte slice to string.
func String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Bytes converts string to byte slice.
func Bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
