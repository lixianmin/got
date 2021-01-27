package convert

import "strconv"

/********************************************************************
created:    2020-01-27
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func AppendInt(b []byte, v interface{}, base int) []byte {
	switch v := v.(type) {
	case int:
		return strconv.AppendInt(b, int64(v), base)
	case int8:
		return strconv.AppendInt(b, int64(v), base)
	case int16:
		return strconv.AppendInt(b, int64(v), base)
	case int32:
		return strconv.AppendInt(b, int64(v), base)
	case int64:
		return strconv.AppendInt(b, v, base)
	case uint:
		return strconv.AppendUint(b, uint64(v), base)
	case uint8:
		return strconv.AppendUint(b, uint64(v), base)
	case uint16:
		return strconv.AppendUint(b, uint64(v), base)
	case uint32:
		return strconv.AppendUint(b, uint64(v), base)
	case uint64:
		return strconv.AppendUint(b, v, base)
	}

	return b
}
