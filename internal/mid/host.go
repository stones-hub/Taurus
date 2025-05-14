package mid

import (
	"Taurus/pkg/httpx"
	"Taurus/pkg/util"
	"net/http"
)

func HostMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ips := util.GetRemoteIP(r)

		// 检查主机是否在允许列表中
		allowedHosts := getAllowedHosts()
		allowed := false

		// 遍历ips是不是在allowedHosts中
		for _, ip := range ips {
			for _, allowedHost := range allowedHosts {
				if util.IsIPAllowed(ip, []string{allowedHost}) {
					allowed = true
					break
				}
			}
		}

		if !allowed {
			httpx.SendResponse(w, http.StatusForbidden, "访问被拒绝：未授权的主机", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getAllowedHosts() []string {
	// 默认允许的主机列表
	return []string{
		"14.18.194.140",
		"14.18.194.128",
		"127.0.0.1",
		"localhost",
		"192.168.0.0/16",
	}
}

/*
- 192.0.0.0/8 192.168.0.0/16 192.168.1.0/24
- /24：子网掩码为255.255.255.0，表示前24位是网络部分，后8位是主机部分。
- /16：子网掩码为255.255.0.0，表示前16位是网络部分，后16位是主机部分。
- /8：子网掩码为255.0.0.0，表示前8位是网络部分，后24位是主机部分。
*/
