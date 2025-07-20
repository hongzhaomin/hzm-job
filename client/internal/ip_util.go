package internal

import (
	"fmt"
	"github.com/hongzhaomin/hzm-job/client/internal/prop"
	"github.com/hongzhaomin/hzm-job/core/ezconfig"
	"net"
	"os"
)

const (
	anyHostValue   = "0.0.0.0"
	localhostValue = "127.0.0.1"
)

// GetHostIp 自动检测环境并返回对应IP
func GetHostIp() string {
	// 优先从配置文件取
	clientConfig := ezconfig.Get[*prop.HzmJobConfigBean]()
	if clientConfig.Localhost != "" {
		return clientConfig.Localhost
	}

	// 1. 尝试获取本地IP（物理机或容器内都适用）
	if ip, err := getLocalIP(); err == nil && isValidIp(ip) {
		return ip
	}

	// 2. 特殊处理Docker默认网桥情况（可选）
	if hostname, err := os.Hostname(); err == nil {
		if addr, err := net.LookupIP(hostname); err == nil {
			for _, ip := range addr {
				if ip.To4() != nil && isValidIp(ip.String()) {
					return ip.String()
				}
			}
		}
	}

	// 3. 回退到默认网关IP（适用于某些Docker网络模式）
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err == nil {
		defer conn.Close()
		addr := conn.LocalAddr().(*net.UDPAddr).IP.String()
		if isValidIp(addr) {
			return addr
		}
	}

	return localhostValue // 最终回退
}

// 获取本机有效IP地址（优先非回环IPv4）
func getLocalIP() (string, error) {
	ifts, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("no valid IP found")
	}
	flags := net.FlagUp | net.FlagRunning
	for _, ift := range ifts {
		if ift.Flags&flags != flags {
			continue
		}
		addrs, err := ift.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					return ipNet.IP.String(), nil
				}
			}
			if ipNet, ok := addr.(*net.IPAddr); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					return ipNet.IP.String(), nil
				}
			}
		}
	}
	return "", fmt.Errorf("no valid IP found")
}

func isValidIp(address string) bool {
	ip := net.ParseIP(address)
	if ip == nil {
		return false
	}

	if ip.To4() != nil {
		return address != anyHostValue && address != localhostValue
	}
	return true
}
