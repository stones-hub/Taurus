package util

import (
	"bufio"
	"fmt"
	"os"
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
