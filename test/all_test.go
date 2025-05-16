package test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

const (
	API_URL       = "http://127.0.0.1:9001"
	AUTHORIZATION = "Bearer test"
)

func TestAPIEndpoints(t *testing.T) {
	tests := []struct {
		url    string
		method string
		body   string
	}{
		{url: "/v1/api/", method: "GET", body: ``},
	}

	for _, test := range tests {
		req, err := http.NewRequest(test.method, API_URL+test.url, strings.NewReader(test.body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		req.Header.Set("Authorization", AUTHORIZATION)
		req.Header.Set("Content-Type", "application/json")

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
	}
}

// go test -v -run test/all_test.go or go test
