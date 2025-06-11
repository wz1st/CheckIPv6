#!/bin/sh

if [ -z "$1" ]; then
  echo "使用方法: $0 <IP地址>"
  exit 1
fi

# 传入的IP地址
IP="$1"

cp -r CheckIPv6 /usr/local/bin/
chmod +x /usr/local/bin/CheckIPv6

cp -r check-node.service /etc/systemd/system/

sed -i "s/192\.168\.31\.2/$IP/g" /etc/systemd/system/check-node.service

systemctl daemon-reload
systemctl enable check-node.service
systemctl start check-node.service