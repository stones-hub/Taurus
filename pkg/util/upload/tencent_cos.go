package upload

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type TencentCOS struct {
	Bucket     string
	Region     string
	SecretID   string
	SecretKey  string
	PathPrefix string
	BaseURL    string
}

// UploadFile upload file to COS
func (tencentCOS *TencentCOS) UploadFile(file *multipart.FileHeader) (string, string, error) {
	client := NewClient(tencentCOS)
	f, openError := file.Open()
	if openError != nil {
		return "", "", errors.New("function file.Open() failed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)

	_, err := client.Object.Put(context.Background(), tencentCOS.PathPrefix+"/"+fileKey, f, nil)
	if err != nil {
		panic(err)
	}
	return tencentCOS.BaseURL + "/" + tencentCOS.PathPrefix + "/" + fileKey, fileKey, nil
}

// DeleteFile delete file form COS
func (tencentCOS *TencentCOS) DeleteFile(key string) error {
	client := NewClient(tencentCOS)
	name := tencentCOS.PathPrefix + "/" + key
	_, err := client.Object.Delete(context.Background(), name)
	if err != nil {
		return errors.New("function bucketManager.Delete() failed, err:" + err.Error())
	}
	return nil
}

// NewClient init COS client
func NewClient(tencentCOS *TencentCOS) *cos.Client {
	urlStr, _ := url.Parse("https://" + tencentCOS.Bucket + ".cos." + tencentCOS.Region + ".myqcloud.com")
	baseURL := &cos.BaseURL{BucketURL: urlStr}
	client := cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  tencentCOS.SecretID,
			SecretKey: tencentCOS.SecretKey,
		},
	})
	return client
}
