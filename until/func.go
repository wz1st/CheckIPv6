package until

import (
	"CheckIPv6/global"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func CheckPort(port string) bool {
	fmt.Println("检查端口65535占用...")
	// 尝试监听给定端口
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Println("端口" + port + "被占用...")
		return false // 端口被占用
	}
	listener.Close() // 关闭监听器
	return true      // 端口未被占用
}

func CheckIp() bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("http://test.ipw.cn")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	if isIPv6(strings.TrimSpace(string(bodyBytes))) {
		return true
	}
	return false
}

func isIPv6(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() == nil && strings.Contains(ip, ":")
}

func GetLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		// 排除虚拟网卡和不活动的网卡
		if iface.Flags&net.FlagUp == 0 {
			continue // 网卡未启用
		}
		name := strings.ToLower(iface.Name)
		if strings.Contains(name, "vmware") ||
			strings.Contains(name, "docker") ||
			strings.Contains(name, "vbox") ||
			strings.Contains(name, "virtual") ||
			strings.Contains(name, "loopback") ||
			strings.Contains(name, "tun") ||
			strings.Contains(name, "tap") ||
			strings.Contains(name, "bridge") ||
			strings.Contains(name, "vnic") ||
			strings.Contains(name, "qmi") ||
			strings.Contains(name, "usb") ||
			strings.Contains(name, "wwan") ||
			strings.Contains(name, "wlan") ||
			strings.Contains(name, "hyper-v") ||
			strings.Contains(name, "qemu") {
			continue // 跳过虚拟网卡
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipNet.IP
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // 只处理IPv4
			}
			if isPrivateIP(ip) {
				return ip.String(), nil
			}
		}
	}

	return "", fmt.Errorf("找不到局域网IP")
}

// 判断是否是私有IP
func isPrivateIP(ip net.IP) bool {
	privateBlocks := []string{"10.", "172.", "192.168."}
	ipStr := ip.String()
	for _, block := range privateBlocks {
		if len(ipStr) >= len(block) && ipStr[:len(block)] == block {
			return true
		}
	}
	return false
}

func NodeCheck(ip string) bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("http://%s:65535/check", ip))
	if err != nil {
		fmt.Println("节点不可用")
		global.HandleTimeout(ip)
		return false
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		global.HandleTimeout(ip)
		return false
	}
	global.ResetTimeout(ip)
	return strings.TrimSpace(string(bodyBytes)) == "1"
}

func Add(ip string, server string) error {
	targetURL := fmt.Sprintf("http://%s:65535/manage", server)

	// 表单数据
	data := url.Values{}
	data.Set("ip", ip)

	// 发送POST请求
	resp, err := http.PostForm(targetURL, data)
	if err != nil {
		fmt.Println("请求失败:", err)
		return err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return err
	}

	if string(body) == "1" {
		return nil
	}
	return errors.New("添加失败")
}
