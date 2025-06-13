package controller

import (
	"Taurus/pkg/contextx"
	"Taurus/pkg/db"
	"Taurus/pkg/httpx"
	"net/http"

	"github.com/google/wire"
)

// Author: yelei
// Email: 61647649@qq.com
// Date: 2025-06-13

/*
TestURL : http://127.0.0.1:9080/v1/api/?age=20&email=yelei@3k.com&id=1&phone=13631375979&name=yelei
Note: 测试validate中间件，验证请求来源的结构体是否符合ValidateRequest结构体
*/

type ValidateCtrl struct {
}

type ValidateRequest struct {
	Id    int    `json:"id" validate:"required,numeric"`
	Name  string `json:"name" validate:"required,min=2,max=10"`
	Age   int    `json:"age" validate:"required,numeric,min=18,max=100"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required,len=11"`
}

var ValidateCtrlSet = wire.NewSet(wire.Struct(new(ValidateCtrl), "*"))

func (validateCtrl *ValidateCtrl) TestValidateMiddleware(w http.ResponseWriter, r *http.Request) {

	// 验证请求是否符合ValidateRequest结构体
	validateRequest, ok := contextx.GetValidateRequest(r.Context())
	if !ok {
		httpx.SendResponse(w, http.StatusBadRequest, "validate request failed", nil)
		return
	}

	// 将请求转换为ValidateRequest结构体
	req, ok := validateRequest.(*ValidateRequest)
	if !ok {
		httpx.SendResponse(w, http.StatusBadRequest, "validate request failed", nil)
		return
	}

	// 测试查询数据库
	var question Question
	if err := db.Find("kf_ai", &question, "id = ?", 1); err != nil {
		httpx.SendResponse(w, http.StatusInternalServerError, "find question failed", nil)
		return
	}

	httpx.SendResponse(w, http.StatusOK, struct {
		ValidateRequest *ValidateRequest `json:"validate_request"`
		Request         interface{}      `json:"request"`
		Question        Question         `json:"question"`
	}{
		ValidateRequest: req,
		Request:         r.URL.Query(),
		Question:        question,
	}, nil)
}

type Question struct {
	Id          int    `json:"id" gorm:"column:id"`
	WorkOrderId string `json:"work_order_id" gorm:"column:work_order_id"`
}

func (question *Question) TableName() string {
	return "questions"
}
