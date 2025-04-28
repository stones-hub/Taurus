package util

import "net"

func GetLocalIPs() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	var ips []string
	for _, addr := range addrs {
		// 检查是否为 IP net.Addr 类型
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		// 获取 IP 地址
		ip := ipNet.IP

		// 排除回环地址
		if ip.IsLoopback() {
			continue
		}

		// 添加 IP 地址
		if ip.To4() != nil || ip.To16() != nil {
			ips = append(ips, ip.String())
		}
	}

	return ips, nil
}
