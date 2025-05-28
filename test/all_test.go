package test

import (
	"Taurus/pkg/redisx"
	"Taurus/pkg/util"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	API_URL       = "http://127.0.0.1:9080"
	AUTHORIZATION = "Bearer 123456"
	USER_AGENT    = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)"
)

type validateRequest struct {
	Id    int    `json:"id" validate:"required,numeric"`
	Name  string `json:"name" validate:"required,min=2,max=10"`
	Age   int    `json:"age" validate:"required,numeric,min=18,max=100"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required,len=11"`
}

func TestAPIEndpoints(t *testing.T) {
	// 初始化redis
	initRedis()

	// 模拟登陆成功， 签发token， 注意redis
	token := loginSuccess(USER_AGENT, 1001, "user-name-test")

	r := validateRequest{
		Id:    1,
		Name:  "test",
		Age:   20,
		Email: "test@test.com",
		Phone: "13800138000",
	}

	b, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	tests := []struct {
		url    string
		method string
		body   string
	}{
		{url: "/v1/api/", method: "POST", body: string(b)},
		{url: "/trace", method: "GET", body: ``},
		{url: "/mid", method: "POST", body: string(b)},
		{url: "/consul/services", method: "GET", body: ``},
		{url: "/health", method: "GET", body: ``},
		{url: "/static/index.html", method: "GET", body: ``},
		{url: "/", method: "GET", body: ``},
	}

	for _, test := range tests {
		req, err := http.NewRequest(test.method, API_URL+test.url, strings.NewReader(test.body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		req.Header.Set("Authorization", AUTHORIZATION)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", USER_AGENT)

		// 添加jwt的token
		req.Header.Set("token", token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", resp.StatusCode)
		}

		t.Logf("Response body: %s", string(body))

		/*
			var response struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Data    interface{} `json:"data"`
			}

			err = json.Unmarshal(body, &response)
			if err != nil {
				t.Fatalf("Failed to parse JSON response: %v", err)
			}

			if response.Code != 200 {
				t.Errorf("Expected code 200, got %d", response.Code)
			}
		*/
	}
}

// 模拟登陆成功， 签发token
func loginSuccess(ua string, userid uint, username string) string {
	token, err := util.GenerateToken(userid, username)
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}

	m := map[string]string{ua: token}

	err = redisx.Redis.HSet(context.Background(), strconv.FormatUint(uint64(userid), 10), m)
	if err != nil {
		log.Fatalf("Failed to set token to redis: %v", err)
	}

	if err != nil {
		log.Fatalf("Failed to set token to redis: %v", err)
	}

	return token
}

// 初始化redis
func initRedis() {
	redisx.InitRedis(redisx.RedisConfig{
		Addrs:        []string{"127.0.0.1:6379"},
		Password:     "",
		DB:           0,
		PoolSize:     100,
		MinIdleConns: 10,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		MaxRetries:   3,
	})
}

// go test -v -run test/all_test.go or go test
