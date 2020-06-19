package convert

import "encoding/json"

/********************************************************************
created:    2019-06-19
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var _toJson = json.Marshal
var _fromJson = json.Unmarshal

func InitJson(toJson func(v interface{}) ([]byte, error), fromJson func(data []byte, v interface{}) error) {
	if toJson == nil {
		panic("toJson is nil")
	}

	if fromJson == nil {
		panic("fromJson is nil.")
	}

	_toJson = toJson
	_fromJson = fromJson
}

func ToJson(v interface{}) ([]byte, error) {
	return _toJson(v)
}

func FromJson(data []byte, v interface{}) error {
	return _fromJson(data, v)
}
