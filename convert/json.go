package convert

import "encoding/json"

/********************************************************************
created:    2019-06-19
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var _toJson = json.Marshal
var _fromJson = json.Unmarshal

func InitJson(marshal func(v any) ([]byte, error), unmarshal func(data []byte, v any) error) {
	if marshal == nil {
		panic("marshal is nil")
	}

	if unmarshal == nil {
		panic("unmarshal is nil.")
	}

	_toJson = marshal
	_fromJson = unmarshal
}

func ToJsonE(v any) ([]byte, error) {
	return _toJson(v)
}

func ToJson(v any) []byte {
	var result, _ = _toJson(v)
	return result
}

func ToJsonS(v any) string {
	var result, _ = _toJson(v)
	return String(result)
}

func FromJsonE(data []byte, v any) error {
	return _fromJson(data, v)
}

func FromJson(data []byte, v any) {
	_ = _fromJson(data, v)
}

func FromJsonS(data string, v any) {
	_ = _fromJson(Bytes(data), v)
}
