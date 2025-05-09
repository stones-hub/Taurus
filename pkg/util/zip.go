package util

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 解压
func Unzip(zipFile string, destDir string) ([]string, error) {
	zipReader, err := zip.OpenReader(zipFile)
	var paths []string
	if err != nil {
		return []string{}, err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		if strings.Contains(f.Name, "..") {
			return []string{}, fmt.Errorf("%s 文件名不合法", f.Name)
		}
		fpath := filepath.Join(destDir, f.Name)
		paths = append(paths, fpath)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return []string{}, err
			}

			inFile, err := f.Open()
			if err != nil {
				return []string{}, err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return []string{}, err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return []string{}, err
			}
		}
	}
	return paths, nil
}

// 对字符串做压缩
func CompressGzip(data string) (string, error) {
	var (
		buf bytes.Buffer
		err error
		w   *gzip.Writer
	)

	w = gzip.NewWriter(&buf)

	if _, err = w.Write([]byte(data)); err != nil {
		return "", err
	}

	w.Close()

	return buf.String(), nil
}

// 字符串解压
func DecompressGzip(data string) (string, error) {
	var (
		buf              bytes.Buffer
		err              error
		r                *gzip.Reader
		decompressedData []byte
	)

	buf.Write([]byte(data))
	if r, err = gzip.NewReader(&buf); err != nil {
		return "", err
	}

	if decompressedData, err = io.ReadAll(r); err != nil {
		return "", err
	}

	return string(decompressedData), nil
}
