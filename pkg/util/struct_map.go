package util

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

// 将结构体转换为map
// obj 结构体或结构体指针
func StructToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	objValue := reflect.ValueOf(obj)

	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}

	if objValue.Kind() != reflect.Struct {
		fmt.Println("输入不是结构体类型")
		return result
	}

	typeOfObj := objValue.Type()
	for i := 0; i < objValue.NumField(); i++ {
		field := objValue.Field(i)
		fieldType := typeOfObj.Field(i)
		tag := fieldType.Tag.Get("json")
		if tag == "" {
			tag = fieldType.Name
		}
		result[tag] = field.Interface()
	}
	return result
}

// 将结构体转换为map[string]string
// obj 结构体或结构体指针
func StructToStringMap(obj interface{}) map[string]string {
	result := make(map[string]string)
	objValue := reflect.ValueOf(obj)

	if objValue.Kind() == reflect.Ptr {
		objValue = objValue.Elem()
	}

	if objValue.Kind() != reflect.Struct {
		log.Printf("输入不是结构体类型: %v \n", objValue.Kind())
		return result
	}

	typeOfObj := objValue.Type()
	for i := 0; i < objValue.NumField(); i++ {
		field := objValue.Field(i)
		fieldType := typeOfObj.Field(i)
		tag := fieldType.Tag.Get("json")
		if tag == "" {
			tag = fieldType.Name
		}
		var valueStr string
		switch field.Kind() {
		case reflect.String:
			valueStr = field.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			valueStr = strconv.FormatInt(field.Int(), 10)
		case reflect.Float32, reflect.Float64:
			valueStr = strconv.FormatFloat(field.Float(), 'f', -1, 64)
		default:
			valueStr = fmt.Sprintf("%v", field.Interface())
		}
		result[tag] = valueStr
	}
	return result
}

// MapToStruct 将 map 转换为结构体
// m map[string]interface{} key 是string，value 是任意类型的map
// s interface{} 结构体指针 &sample
func MapToStruct(m map[string]interface{}, s interface{}) error {
	// 将 map 转换为 JSON
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	// 将 JSON 转换为结构体
	return json.Unmarshal(jsonData, s)
}

// map -> struct
// map -> json -> struct

// struct -> map
// 反射最方便
