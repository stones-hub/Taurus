package controller

import (
	"Taurus/pkg/consul"
	"Taurus/pkg/httpx"
	"log"
	"net/http"

	"github.com/google/wire"
)

type ConsulCtrl struct {
}

var ConsulCtrlSet = wire.NewSet(wire.Struct(new(ConsulCtrl), "*"))

func (c *ConsulCtrl) TestConsul(w http.ResponseWriter, r *http.Request) {
	// 获取consul注册的服务
	services, err := consul.Client.GetServices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 打印consul注册的服务
	log.Println("consul注册的服务 -> ", services)

	// 发现服务
	service, err := consul.Client.Discover("taurus-api-gateway")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 打印发现的服务
	log.Println("服务详情: ",
		service.Node.ID,
		service.Service.ID,
		service.Service.Tags,
		service.Service.Address,
		service.Service.Port,
		service.Service.Meta,
		service.Service.Kind,
		service.Service.TaggedAddresses,
		service.Service.Weights,
		service.Service.EnableTagOverride,
		service.Service.CreateIndex,
		service.Service.ModifyIndex,
		service.Checks,
	)

	// 注销consul注册的服务
	err = consul.Client.Deregister(service.Service.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 打印注销后的consul注册的服务
	services, err = consul.Client.GetServices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 打印consul注册的服务
	httpx.SendResponse(w, http.StatusOK, services, nil)
}
