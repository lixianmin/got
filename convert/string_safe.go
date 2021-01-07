// +build appengine

package convert

/********************************************************************
created:    2021-01-06
author:     lixianmin

this file is extracted from go-redis/v8/internal/safe.go
*********************************************************************/

func String(b []byte) string {
	return string(b)
}

func Bytes(s string) []byte {
	return []byte(s)
}
