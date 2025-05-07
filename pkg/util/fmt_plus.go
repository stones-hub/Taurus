package util

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: StructToMap
//@description: 利用反射将结构体转化为map
//@param: obj interface{}
//@return: map[string]interface{}

func StructToMapByeTagmapstructure(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	data := make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		if obj1.Field(i).Tag.Get("mapstructure") != "" {
			data[obj1.Field(i).Tag.Get("mapstructure")] = obj2.Field(i).Interface()
		} else {
			data[obj1.Field(i).Name] = obj2.Field(i).Interface()
		}
	}
	return data
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: ArrayToString
//@description: 将数组格式化为字符串
//@param: array []interface{}
//@return: string

func ArrayToString(array []interface{}) string {
	return strings.Replace(strings.Trim(fmt.Sprint(array), "[]"), " ", ",", -1)
}

func Pointer[T any](in T) (out *T) {
	return &in
}

func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// MaheHump 将字符串转换为驼峰命名，支持自定义分隔符
func MaheHump(s, delimiter string) string {
	// Replace the custom delimiter with a space
	s = strings.ReplaceAll(s, delimiter, " ")
	words := strings.Fields(s)
	c := cases.Title(language.Und)

	for i := 1; i < len(words); i++ {
		words[i] = c.String(words[i])
	}

	return strings.Join(words, "")
}

// 随机字符串
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:',.<>?/~`")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[RandomInt(0, len(letters))]
	}
	return string(b)
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// FilterInvisibleChars removes invisible characters from a string.
// ASCII码，通常小于等于32或者大于等于127的都属于不可见字符, 比如空格, 换行符, 制表符等，过滤掉
func FilterInvisibleChars(s string) string {
	resRunes := []rune{}
	for _, r := range s {
		// ASCII码，通常小于等于32或者大于等于127的都属于不可见字符
		if r > 32 && r < 127 {
			resRunes = append(resRunes, r)
		}
	}
	return string(resRunes)
}
