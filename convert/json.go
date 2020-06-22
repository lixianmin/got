package convert

import "encoding/json"

/********************************************************************
created:    2019-06-19
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var _toJson = json.Marshal
var _fromJson = json.Unmarshal

func InitJson(marshal func(v interface{}) ([]byte, error), unmarshal func(data []byte, v interface{}) error) {
	if marshal == nil {
		panic("marshal is nil")
	}

	if unmarshal == nil {
		panic("unmarshal is nil.")
	}

	_toJson = marshal
	_fromJson = unmarshal
}

func ToJson(v interface{}) ([]byte, error) {
	return _toJson(v)
}

func FromJson(data []byte, v interface{}) error {
	return _fromJson(data, v)
}
