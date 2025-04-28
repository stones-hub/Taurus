package util

import "strconv"

func ParseStringSliceToUint64(s []string) []uint64 {
	iv := make([]uint64, len(s))

	for i, v := range s {
		// 以10进制的方式解析v， 最后保存为64 uint
		iv[i], _ = strconv.ParseUint(v, 10, 64)
	}

	/**
	// 将s 字符串用 base 进制转成 bitSize位的int类型
	strconv.ParseInt(s string, base int, bitSize int)
	// 将s 字符串用 base 进制转成 bitSize位的uint类型
	strconv.ParseUint(s string, base int , bitSize int)
	*/
	return iv
}
