package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// 通用的HTTP请求函数
func doHttpRequest(method, url string, payload interface{}, headers map[string]string) ([]byte, error) {
	var (
		err         error
		request     *http.Request
		response    *http.Response
		body        []byte
		jsonPayload []byte
	)

	// 如果payload是map或字符串类型，进行处理
	if payload != nil {
		switch p := payload.(type) {
		case []byte:
			jsonPayload = p
		case string:
			jsonPayload = []byte(p) // 将字符串转换为字节数组
		default:
			if jsonPayload, err = json.Marshal(payload); err != nil {
				return nil, err
			}
		}
	}

	// 创建HTTP请求
	if request, err = http.NewRequest(method, url, bytes.NewBuffer(jsonPayload)); err != nil {
		return nil, err
	}

	// 设置请求头
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}

	// 发送请求
	if response, err = http.DefaultClient.Do(request); err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// 读取响应体
	if body, err = io.ReadAll(response.Body); err != nil {
		return nil, err
	}

	return body, nil
}

// POST请求的封装
func HttpPost(url string, payload interface{}, headers map[string]string) ([]byte, error) {
	return doHttpRequest("POST", url, payload, headers)
}

// GET请求的封装
func HttpGet(url string, headers map[string]string) ([]byte, error) {
	return doHttpRequest("GET", url, nil, headers)
}

/*
HttpClient请求 封装
*/
func HttpRequest(URL string, method string, headers map[string]string, params map[string]string, data any) (*http.Response, error) {
	var (
		err      error
		u        *url.URL
		query    url.Values
		body     = &bytes.Buffer{} // 设置body数据
		dataJson []byte
		req      *http.Request
		resp     *http.Response
	)
	// 创建URL
	u, err = url.Parse(URL)
	if err != nil {
		return nil, err
	}

	// 添加查询参数
	query = u.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	u.RawQuery = query.Encode()

	// 将数据编码为JSON
	if data != nil {
		dataJson, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
		//  fmt.Println("http send data:", string(bodyData))
		body = bytes.NewBuffer(dataJson)
	}

	// 创建请求
	req, err = http.NewRequest(method, u.String(), body)

	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// 发送请求
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// 返回响应，让调用者处理
	return resp, nil
}

/*
获取http请求返回的数据
*/
func ReadResponse(res *http.Response) ([]byte, error) {
	return io.ReadAll(res.Body)
}
