package util

import (
	"fmt"
	"testing"
)

func TestParseStringSliceToUint64(t *testing.T) {
	iv := ParseStringSliceToUint64([]string{"中国", "人民", "站起来了"})
	fmt.Println("---------- uint slice --------------------")
	fmt.Println(iv)
}
