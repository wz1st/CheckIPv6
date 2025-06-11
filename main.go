package main

import (
	"CheckIPv6/boot"
	"CheckIPv6/router"
	"CheckIPv6/until"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	manage := flag.Bool("m", false, "以主节点方式运行")
	server := flag.String("s", "192.168.31.2", "IPv6检测主节点IP地址")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "使用方式: %s [options]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}
	if *manage && *server != "192.168.31.2" {
		fmt.Println("错误: -m 和 -s 不能同时设置")
		flag.Usage()
		os.Exit(1)
	}
	if !until.CheckPort("65535") {
		return
	}
	if boot.Init(*manage, *server) != nil {
		return
	}
	fmt.Println("启动接口...")
	router := router.InitRouter(*manage)
	fmt.Println("启动成功")
	router.Run(":" + "65535")
}
