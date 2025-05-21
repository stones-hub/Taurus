package consuls

import (
	"Taurus/pkg/consul"
	"log"

	"github.com/hashicorp/consul/api"
)

// 实现ttlupdate接口
type DefaultTTLUpdater struct {
}

// 更新TTL
func (u *DefaultTTLUpdater) Update(c *consul.ConsulClient, checkID string) error {
	log.Printf("更新TTL..")
	c.UpdateTTL(checkID, api.HealthPassing, "TTL更新")
	return nil
}
