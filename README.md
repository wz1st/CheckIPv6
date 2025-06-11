### why
OpenWRT(或基于OP做的路由器)做二级路由经常性的掉PD，做了这个check脚本，可以自动重启wan6，避免掉线

### how to use
#### 单路由器检测
1. 修改脚本`check_ipv6.sh`中的`wan6`为你的OpenWRT 的wan6接口名
2. 上传`op-install.sh` `check_ipv6.sh`到OP上，执行`sh op-install.sh`即可。

#### 路由器+内网节点检测
1. 修改脚本`check_ipv6.sh`中的`wan6`为你的OpenWRT 的wan6接口名
2. 主节点安装：上传`install-manage.sh` `CheckIPv6` `check-server.service`到OP可以ping通的Linux上，执行`sh install-manage.sh`即可。
3. 子节点安装：上传`install-node.sh` `CheckIPv6` `check-node.service`到内网节点上，执行`sh install-node.sh xxx`即可。
4. 上传`op-install.sh` `check_ipv6.sh`到OP上，执行`sh op-install.sh xxx`即可。
##### 参数说明
>上面的`xxx`是主节点的IP地址，仅支持一个主节点。子节点为可选项，不安装则默认不检测子节点。