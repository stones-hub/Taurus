package util

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// 按行写文件
func WriteLine(filename string, b []byte) error {
	var (
		err    error
		fd     *os.File
		writer *bufio.Writer
	)
	fd, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()
	writer = bufio.NewWriter(fd)

	_, err = writer.WriteString(string(b) + "\n")
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}

// 一次读取文件所有数据
func ReadAll(filename string) (string, error) {

	var (
		err     error
		fd      *os.File
		scanner *bufio.Scanner
		content string
	)

	fd, err = os.OpenFile(filename, os.O_RDONLY|os.O_APPEND, os.ModePerm)

	if err != nil {
		return "", err
	}
	defer fd.Close()

	scanner = bufio.NewScanner(fd)

	for scanner.Scan() {
		line := scanner.Text()
		content += line
	}

	if err = scanner.Err(); err != nil {
		return "", err
	}

	if len(content) == 0 {
		return "", fmt.Errorf("文件为空")
	}

	return content, nil
}

// 按行读取文件，返回[]interface{}
func ReadLine(filename string) ([]interface{}, error) {
	var (
		err     error
		fd      *os.File
		scanner *bufio.Scanner
		content []interface{}
	)

	fd, err = os.OpenFile(filename, os.O_RDONLY|os.O_APPEND, os.ModePerm)

	if err != nil {
		return nil, err
	}
	defer fd.Close()

	scanner = bufio.NewScanner(fd)

	for scanner.Scan() {
		line := scanner.Text()
		content = append(content, line)
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return content, nil
}

// FetchAllDir 递归遍历文件目录
func FetchAllDir(path string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(path, func(filepath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, filepath)
		}
		return nil
	})

	return files, err
}

// WalkDir 递归遍历文件目录, 含深度
func WalkDir(path string, maxDepth int) ([]string, error) {

	var files []string

	// 计算根目录的路径深度
	rootDepth := len(strings.Split(path, string(os.PathSeparator)))

	err := filepath.WalkDir(path, func(currentPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 计算当前路径的深度
		currentDepth := len(strings.Split(currentPath, string(os.PathSeparator))) - rootDepth

		// 检查是否超出最大深度
		if currentDepth > maxDepth {
			return filepath.SkipDir // 跳过此目录及其子目录
		}

		// 如果是文件，则添加到结果切片
		if !d.IsDir() {
			files = append(files, currentPath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// CheckPath 检查路径是否是文件or目录
func CheckPath(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	if info.IsDir() {
		return true, nil
	}

	return false, nil
}

// GetCurrentPath 读取当前文件所在目录
func GetCurrentPath() (string, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("failed to get caller info: %v", ok)
	}
	dir := filepath.Dir(file)
	return dir, nil
}

// PathExists 文件目录是否存在
func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return true, nil
		}
		return false, errors.New("存在同名文件")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CreateDir 批量创建文件夹
func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := PathExists(v)
		if err != nil {
			return err
		}
		if !exist {
			if err := os.MkdirAll(v, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return err
}

// FileMove 文件移动
//@param: src string, dst string(src: 源位置,绝对路径or相对路径, dst: 目标位置,绝对路径or相对路径,必须为文件夹)
//@return: err error

func FileMove(src string, dst string) (err error) {
	if dst == "" {
		return nil
	}
	src, err = filepath.Abs(src)
	if err != nil {
		return err
	}
	dst, err = filepath.Abs(dst)
	if err != nil {
		return err
	}
	revoke := false
	dir := filepath.Dir(dst)
Redirect:
	_, err = os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0o755)
		if err != nil {
			return err
		}
		if !revoke {
			revoke = true
			goto Redirect
		}
	}
	return os.Rename(src, dst)
}

//@description: 去除结构体空格
//@param: target interface (target: 目标结构体,传入必须是指针类型)

func TrimSpace(target interface{}) {
	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Ptr {
		return
	}
	t = t.Elem()
	v := reflect.ValueOf(target).Elem()
	for i := 0; i < t.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.String:
			v.Field(i).SetString(strings.TrimSpace(v.Field(i).String()))
		}
	}
}

// FileExist 判断文件是否存在
func FileExist(path string) bool {
	fi, err := os.Lstat(path)
	if err == nil {
		return !fi.IsDir()
	}
	return !os.IsNotExist(err)
}
