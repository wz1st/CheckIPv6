package global

import (
	"sync"
)

// IPStore 全局保存的IP地址集合（示例）
var (
	IPStore = make(map[string]bool) // 用 map 存储，去重方便
	IPMux   sync.RWMutex
)

func GetAllIPs() []string {
	IPMux.RLock()
	defer IPMux.RUnlock()

	ips := make([]string, 0, len(IPStore))
	for ip := range IPStore {
		ips = append(ips, ip)
	}
	return ips
}

func AddIP(ip string) {
	IPMux.Lock()
	defer IPMux.Unlock()
	IPStore[ip] = true
}
