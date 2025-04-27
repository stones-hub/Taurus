package util

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
)

// 将数据写入CSV文件, data数据是没有header头的
// filename csv文件地址
// data 待写入的数据
// csv头数据，针对需要固定头位置的场景
func GenCSV(filename string, data []map[string]string, headers []string) error {
	var (
		fd     *os.File
		err    error
		writer *csv.Writer
	)

	if fd, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm); err != nil {
		log.Printf("生成csv文件失败: %s \n", err.Error())
		return err
	}

	defer fd.Close()

	writer = csv.NewWriter(fd)

	defer writer.Flush()

	if len(data) == 0 {
		return fmt.Errorf("数据为空")
	}

	if len(headers) == 0 {
		// 创建一个表头大小的slice
		headers = make([]string, 0, len(data[0]))
		for header := range data[0] {
			headers = append(headers, header)
		}
	}

	// 写表头
	if err = writer.Write(headers); err != nil {
		log.Printf("写入表头失败: %s \n", err.Error())
		return err
	}

	// 写数据
	for _, row := range data {

		record := make([]string, 0, len(headers))

		for _, header := range headers {
			record = append(record, row[header])
		}

		// log.Printf("写入csv文件数据： %v \n", record)

		writer.Write(record)
	}

	return nil
}

/*
type Test struct {
	Name string `csv:"name" json:"name"`
	Age  int    `csv:"age" json:"age"`
}
*/
// 将数据从csv中读取出来
// 假如 result : &[]Test{},  result 是slice的指针， slice中存的是结构体，不可以是结构体指针
func ReadCSV(fileanme string, result interface{}) error {
	var (
		fd               *os.File
		err              error
		reader           *csv.Reader
		invisibleHeaders []string
		headers          []string
		resultValue      reflect.Value
	)

	if fd, err = os.OpenFile(fileanme, os.O_RDONLY|os.O_APPEND, os.ModePerm); err != nil {
		return err
	}

	defer fd.Close()
	reader = csv.NewReader(fd)

	// 读取csv文件第一行，默认为表头
	if invisibleHeaders, err = reader.Read(); err != nil {
		return err
	}

	// 由于csv中读取的表头有可能有很多编码且隐藏字符的问题，所以做一次过滤
	for _, v := range invisibleHeaders {
		resRunes := []rune{}
		for _, r := range v {
			// ascii码，通常小于等于32或者大于等于127的都属于不可见字符
			if r > 32 && r < 127 {
				resRunes = append(resRunes, r)
			}
		}
		// 打印ascii编码
		// fmt.Println(v, resRunes)
		headers = append(headers, string(resRunes))
	}

	// 获取result的反射类型
	resultValue = reflect.ValueOf(result)

	// result的类型不是指针 或者 result的值不是slice的话，不可以
	if resultValue.Kind() != reflect.Ptr || resultValue.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("Result必须是slice的指针")
	}

	// 获取result指针指向的slice数组的类型: sliceType = []Test 类型
	sliceType := resultValue.Elem().Type()

	// 获取result指针指向的slice数组的类型的类型值: elementType = Test
	elementType := sliceType.Elem()

	// 根据result指向的slice数组类型创建一个新的slice : []Test
	slice := reflect.MakeSlice(sliceType, 0, 0)

	// 读取数据行
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("读取数据行失败: %v\n", err)
			continue
		}

		// 创建新的结构体实例
		element := reflect.New(elementType).Elem()

		// 遍历结构体的字段
		for i := 0; i < element.NumField(); i++ {

			field := element.Type().Field(i)
			tag := field.Tag.Get("csv")
			if tag == "" {
				continue
			}

			// 找到对应的CSV列索引
			colIndex := -1
			for j, header := range headers {
				if header == tag {
					colIndex = j
					break
				}
			}

			if colIndex == -1 || colIndex >= len(record) {
				continue
			}

			// 设置字段值
			fieldValue := element.Field(i)
			value := record[colIndex]

			switch fieldValue.Kind() {
			case reflect.String:
				fieldValue.SetString(value)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if v, err := strconv.ParseInt(value, 10, 64); err == nil {
					fieldValue.SetInt(v)
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if v, err := strconv.ParseUint(value, 10, 64); err == nil {
					fieldValue.SetUint(v)
				}
			case reflect.Float32, reflect.Float64:
				if v, err := strconv.ParseFloat(value, 64); err == nil {
					fieldValue.SetFloat(v)
				}
			case reflect.Bool:
				if v, err := strconv.ParseBool(value); err == nil {
					fieldValue.SetBool(v)
				}
			case reflect.Slice:
				// 处理 []byte 类型
				if fieldValue.Type().Elem().Kind() == reflect.Uint8 {
					fieldValue.SetBytes([]byte(value))
				}
			default:
				// 处理 time.Time 类型
				if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
					if v, err := time.Parse("2006-01-02 15:04:05", value); err == nil {
						fieldValue.Set(reflect.ValueOf(v))
					}
				} else if fieldValue.Type().Kind() == reflect.Interface {
					// 处理 interface{} 类型
					fieldValue.Set(reflect.ValueOf(value))
				}
			}
		}

		// 将结构体添加到切片中
		slice = reflect.Append(slice, element)
	}

	// 设置结果
	resultValue.Elem().Set(slice)
	return nil
}
