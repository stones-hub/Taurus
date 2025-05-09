package util

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"reflect"
	"strconv"
	"time"
)

// ExcelWriter 类
type ExcelWriter struct {
	file    *os.File
	writer  *csv.Writer
	headers []string
}

// InitExcelWriter 创建一个新的 ExcelWriter 实例
func InitExcelWriter(filename string, headers []string) (*ExcelWriter, error) {
	var (
		err    error
		file   *os.File
		writer *csv.Writer
	)

	// 判断目录是否存在
	file, err = os.Create(filename)
	if err != nil {
		return nil, err
	}

	writer = csv.NewWriter(file)
	//  初始化的时候，先将头写到文件
	if err = writer.Write(headers); err != nil {
		_ = file.Close()
		return nil, err
	}

	return &ExcelWriter{
		file:    file,
		writer:  writer,
		headers: headers,
	}, nil
}

// WriteBatch 写入一批数据, data 数据数组
func (excelWriter *ExcelWriter) WriteBatch(datas []interface{}) error {

	for _, record := range datas {

		// 以excel列头为长度，创建能存储一行数据的slice
		row := make([]string, len(excelWriter.headers))
		// 反射一行数据的结构体对象 {"":}
		rowVal := reflect.ValueOf(record)

		//  Excel头(header)的值和行结构体数据的字段名(KEY)是设置的一样的，否则下面没有办法通过字段名那到数据值
		for i, header := range excelWriter.headers {
			// 根据字段名称，获取字段名称存储的值
			fieldVal := rowVal.FieldByName(header)

			if fieldVal.IsValid() {

				switch fieldVal.Kind() { // 判断字段对应的值的类型

				case reflect.String:
					row[i] = fieldVal.String()

				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					row[i] = strconv.FormatInt(fieldVal.Int(), 10)

				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					row[i] = strconv.FormatUint(fieldVal.Uint(), 10)

				case reflect.Float32, reflect.Float64:
					row[i] = strconv.FormatFloat(fieldVal.Float(), 'f', -1, 64)

				case reflect.Bool:
					row[i] = strconv.FormatBool(fieldVal.Bool())

				case reflect.Struct:
					if fieldVal.Type() == reflect.TypeOf(time.Now()) { // 字段类型是时间结构体time.Time
						// 这个写法很有意思， 先将fieldVal转成interface在转成time.Time
						row[i] = fieldVal.Interface().(time.Time).Format("2006-01-02 15:04:05")
					} else { // 如果不是time.Time结构体类型，统一json成字符串在写入到Excel
						jsonBytes, _ := json.Marshal(fieldVal.Interface())
						row[i] = string(jsonBytes)
					}
				default:
					jsonBytes, _ := json.Marshal(fieldVal.Interface())
					row[i] = string(jsonBytes)
				}
			} else {
				row[i] = ""
			}
		}

		if err := excelWriter.writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}

// Close 关闭 CSV 文件
func (excelWriter *ExcelWriter) Close() error {
	excelWriter.writer.Flush()
	if err := excelWriter.writer.Error(); err != nil {
		return err
	}
	return excelWriter.file.Close()
}
