package util

import (
	"github.com/jinzhu/copier"
)

func Copy(s, ts interface{}) error { // 将一个结构体 copy 到另外一个结构体
	return copier.Copy(ts, s)
}
