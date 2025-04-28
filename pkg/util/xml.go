package util

import (
	"encoding/xml"
	"os"
)

/*
XML文件读取, filePath 文件地址, v 需要映射的结构体指针
*/

func ReadFromXml(filePath string, v interface{}) error {

	fd, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fd.Close()
	decodeXML := xml.NewDecoder(fd)

	if err := decodeXML.Decode(v); err != nil {
		return err
	}

	return nil
}

// XML文件生成， v 需要写入xml文件的结构体指针， filePath 文件地址

func WriteToXml(v interface{}, filePath string) error {
	fd, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fd.Close()
	return xml.NewEncoder(fd).Encode(v)
}
