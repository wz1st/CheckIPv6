package global

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// 文件名常量
const ipFile = "ipstore.json"

// IPStore 保存 IP 地址的集合
var (
	IPStore         = make(map[string]bool) // IP存储
	IPTimeoutCount  = make(map[string]int)  // IP超时次数
	IPMux           sync.RWMutex            // 读写锁
	MaxTimeoutCount = 3                     // 最大超时次数
)

// LoadIPs 从文件加载 IPs
func LoadIPs() error {
	IPMux.Lock()
	defer IPMux.Unlock()
	fmt.Println("加载文件中子节点")

	file, err := os.Open(ipFile)
	if err != nil {
		// 文件不存在可以直接返回
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&IPStore)
}

// SaveIPs 保存 IPs 到文件
func SaveIPs() error {

	file, err := os.Create(ipFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(IPStore)
}

// GetAllIPs 返回所有 IP
func GetAllIPs() []string {
	IPMux.RLock()
	defer IPMux.RUnlock()

	ips := make([]string, 0, len(IPStore))
	for ip := range IPStore {
		ips = append(ips, ip)
	}
	return ips
}

// AddIP 添加一个 IP 并保存到文件
func AddIP(ip string) error {
	IPMux.Lock()
	defer IPMux.Unlock()

	IPStore[ip] = true
	return SaveIPs()
}

func HandleTimeout(ip string) {
	IPMux.Lock()
	defer IPMux.Unlock()

	IPTimeoutCount[ip]++
	if IPTimeoutCount[ip] >= MaxTimeoutCount {
		delete(IPStore, ip)
		delete(IPTimeoutCount, ip)
		fmt.Printf("节点 %s 超时次数达到阈值，已删除\n", ip)
		SaveIPs()
	}
}

// 重置超时次数
func ResetTimeout(ip string) {
	IPMux.Lock()
	defer IPMux.Unlock()

	IPTimeoutCount[ip] = 0
}
