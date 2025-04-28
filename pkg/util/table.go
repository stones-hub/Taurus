package util

import (
	"fmt"
	"strings"
)

// TableData 定义表格数据
type TableData [][]interface{}

// PrintTable 打印表格
func PrintTable(data TableData) {
	// 计算每列的最大宽度
	colWidths := make([]int, len(data[0]))
	for _, row := range data {
		for i, val := range row {
			strVal := formatValue(val)
			if len(strVal) > colWidths[i] {
				colWidths[i] = len(strVal)
			}
		}
	}

	// 打印表格
	for _, row := range data {
		for i, val := range row {
			strVal := formatValue(val)
			// 使用空格填充到最大宽度
			paddedVal := fmt.Sprintf("%-"+fmt.Sprintf("%ds", colWidths[i])+"", strVal)
			fmt.Print(paddedVal)
			if i < len(row)-1 {
				fmt.Print(" | ")
			}
		}
		fmt.Println()
		if len(row) > 0 {
			fmt.Println(strings.Repeat("-", sum(colWidths)+len(row)*2+1))
		}
	}
}

// formatValue 将任意类型的数据转换为字符串
func formatValue(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case []interface{}:
		return fmt.Sprintf("%v", v)
	case map[string]interface{}:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// sum 计算整数切片的总和
func sum(nums []int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}
