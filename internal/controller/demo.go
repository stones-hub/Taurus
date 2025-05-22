package controller

import (
	"Taurus/pkg/contextx"
	"Taurus/pkg/db"
	"Taurus/pkg/httpx"
	"Taurus/pkg/logx"
	"log"
	"net/http"

	"github.com/google/wire"
)

type DemoCtrl struct {
}

type DemoRequest struct {
	Id    int    `json:"id" validate:"required,numeric"`
	Name  string `json:"name" validate:"required,min=2,max=10"`
	Age   int    `json:"age" validate:"required,numeric,min=18,max=100"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required,len=11"`
}

var DemoCtrlSet = wire.NewSet(wire.Struct(new(DemoCtrl), "*"))

func (c *DemoCtrl) Get(w http.ResponseWriter, r *http.Request) {
	data, _ := httpx.ParseJson(r)

	validateRequest, ok := contextx.GetValidateRequest(r.Context())
	if !ok {
		httpx.SendResponse(w, http.StatusBadRequest, "validate request failed", nil)
		return
	}

	req, ok := validateRequest.(*DemoRequest)
	if !ok {
		httpx.SendResponse(w, http.StatusBadRequest, "validate request failed", nil)
		return
	}

	var question Question
	if err := db.Find("kf_ai", &question, "id = ?", 1); err != nil {
		httpx.SendResponse(w, http.StatusInternalServerError, "find question failed", nil)
		return
	}

	log.Printf("question: %v", question)

	logx.Core.Info("default", "question: %v", question)

	httpx.SendResponse(w, http.StatusOK, struct {
		Req  *DemoRequest           `json:"req"`
		Data map[string]interface{} `json:"data"`
	}{
		Req:  req,
		Data: data,
	}, nil)
}

type Question struct {
	Id          int    `json:"id" gorm:"column:id"`
	WorkOrderId string `json:"work_order_id" gorm:"column:work_order_id"`
}

func (question *Question) TableName() string {
	return "questions"
}
