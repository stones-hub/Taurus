// Copyright (c) 2025 Taurus Team. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Author: yelei
// Email: 61647649@qq.com
// Date: 2025-06-13

package controller

import (
	"Taurus/pkg/contextx"
	"Taurus/pkg/httpx"
	"net/http"
	"strconv"
)

// UserRequest 用户请求结构体
type UserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Age      int32  `json:"age" validate:"required,gt=0,lt=150"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

// UpdateUserRequest 更新用户请求结构体
type UpdateUserRequest struct {
	ID    int64  `json:"id" validate:"required,gt=0"`
	Name  string `json:"name" validate:"required,min=2,max=50"`
	Email string `json:"email" validate:"required,email"`
	Age   int32  `json:"age" validate:"required,gt=0,lt=150"`
}

// UserController 用户控制器
type UserController struct{}

// CreateUser 创建用户
func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	// 从上下文中获取已验证的请求数据
	validateReq, ok := contextx.GetValidateRequest(r.Context())
	if !ok {
		httpx.SendResponse(w, http.StatusBadRequest, "无效的请求数据", nil)
		return
	}

	req, ok := validateReq.(*UserRequest)
	if !ok {
		httpx.SendResponse(w, http.StatusBadRequest, "无效的请求数据类型", nil)
		return
	}

	// 模拟创建用户
	user := map[string]string{
		"id":       "1",
		"name":     req.Name,
		"email":    req.Email,
		"age":      strconv.FormatInt(int64(req.Age), 10),
		"password": "******",
	}

	httpx.SendResponse(w, http.StatusOK, "创建用户成功", user)
}

// UpdateUser 更新用户
func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// 从上下文中获取已验证的请求数据
	validateReq, ok := contextx.GetValidateRequest(r.Context())
	if !ok {
		httpx.SendResponse(w, http.StatusBadRequest, "无效的请求数据", nil)
		return
	}

	req, ok := validateReq.(*UpdateUserRequest)
	if !ok {
		httpx.SendResponse(w, http.StatusBadRequest, "无效的请求数据类型", nil)
		return
	}

	// 模拟更新用户
	user := map[string]string{
		"id":    strconv.FormatInt(req.ID, 10),
		"name":  req.Name,
		"email": req.Email,
		"age":   strconv.FormatInt(int64(req.Age), 10),
	}

	httpx.SendResponse(w, http.StatusOK, "更新用户成功", user)
}
