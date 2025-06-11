#!/bin/sh

FLAG_FILE="/tmp/ipv6_prefix_change_flag"
WAN_IFACE="wan_6"

cleanup() {
    echo "0" > "$FLAG_FILE"
}
trap cleanup EXIT

DEFAULT_IP=""
IP="${1:-$DEFAULT_IP}"

FLAG_VALUE=$(cat "$FLAG_FILE" 2>/dev/null)
if [ "$FLAG_VALUE" = "1" ]; then
    logger -t prefix_change "脚本运行中.跳过本次检测"
    exit 0
fi

echo "1" > "$FLAG_FILE"

RESTART_REQUIRED=0

# 检查IPv6默认路由
IPV6_DEFAULT_ROUTE=$(ip -6 route | grep default 2>/dev/null)
if [ -z "$IPV6_DEFAULT_ROUTE" ]; then
    logger -t prefix_change "IPv6 DP丢失."
    RESTART_REQUIRED=1
fi

# 检查公网IP
PUBLIC_IP=$(curl -s --max-time 5 test.ipw.cn)
if echo "$PUBLIC_IP" | grep -qE '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$'; then
    logger -t prefix_change "检测到公网为 IPv4 ($PUBLIC_IP)."
    RESTART_REQUIRED=1
fi

if [ -n "$IP" ]; then
    RES=$(curl -s --max-time 10 "$IP:65535")
    CURL_EXIT_CODE=$?

    if [ $CURL_EXIT_CODE -ne 0 ]; then
        logger -t prefix_change "curl命令执行失败，代码: $CURL_EXIT_CODE."
        RESTART_REQUIRED=1
    elif [ "$RES" = "0" ]; then
        logger -t prefix_change "节点IPv6检测不通过."
        RESTART_REQUIRED=1
    fi
fi

# 如果任何条件触发，则重启接口
if [ "$RESTART_REQUIRED" -eq 1 ]; then
    logger -t prefix_change "Reconnecting $WAN_IFACE interface."
    ifdown "$WAN_IFACE"
    ifup "$WAN_IFACE"
fi