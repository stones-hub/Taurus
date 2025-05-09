package util

import (
	"strconv"
	"strings"
)

// IsIdCard 检查身份证号码是否合法
func IsIdCard(idCard string) bool {

	if len(idCard) != 18 {
		return false
	}

	// 定义校验码映射表
	checkSumMap := map[int]byte{
		0:  '1',
		1:  '0',
		2:  'X',
		3:  '9',
		4:  '8',
		5:  '7',
		6:  '6',
		7:  '5',
		8:  '4',
		9:  '3',
		10: '2',
	}

	// 定义加权因子
	factors := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

	// 将身份证号码转为大写，因为最后一位可能为X
	idCard = strings.ToUpper(idCard)

	// 检查最后一位是否为X
	if idCard[17] == 'X' {
		idCard = idCard[:17] + "a"
	}

	// 计算加权和
	sum := 0
	for i, factor := range factors {
		digit, err := strconv.Atoi(string(idCard[i]))
		if err != nil {
			return false
		}
		sum += digit * factor
	}

	// 获取校验码
	expectedCheckDigit := checkSumMap[sum%11]

	// 检查最后一位是否匹配
	return expectedCheckDigit == idCard[17]

}
