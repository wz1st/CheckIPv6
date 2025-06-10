package router

import (
	"CheckIPv6/global"
	"CheckIPv6/until"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(manage bool) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	router := r.Group("/")
	{
		if manage {
			router.POST("manage", Manage)
			router.GET("check", Check1)
		} else {
			router.GET("check", Check2)
		}
	}
	return r
}

func Manage(c *gin.Context) {
	ip := c.PostForm("ip")

	// 校验 IP 是否为空
	if ip == "" {
		c.String(http.StatusOK, "0")
		return
	}

	// 解析 IP
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil || parsedIP.To4() == nil {
		c.String(http.StatusOK, "0")
		return
	}
	fmt.Println("子节点添加:", ip)

	global.AddIP(ip)

	// 如果通过了校验
	c.String(http.StatusOK, "1")
}

func Check1(c *gin.Context) {
	ipList := global.GetAllIPs()
	a := 0
	if until.CheckIp() {
		fmt.Println("主节点IPv6可用")
		a += 1
	}

	if len(ipList) > 0 {
		fmt.Println("子节点检测")
		for _, ip := range ipList {
			if until.NodeCheck(ip) {
				fmt.Println("子节点:", ip, "IPv6可用")
				a += 1
			} else {
				fmt.Println("子节点:", ip, "IPv6不可用")
			}
		}
		if float64(a)/float64(len(ipList)+1) >= 0.5 {
			fmt.Println("节点IPv6可用大于50%")
			c.String(http.StatusOK, "1")
			return
		}
	}
	c.String(http.StatusOK, "0")
}

func Check2(c *gin.Context) {
	fmt.Println("检测IPv6")
	if until.CheckIp() {
		fmt.Println("子节点IPv6可用")
		c.String(http.StatusOK, "1")
		return
	}
	c.String(http.StatusOK, "0")
}
