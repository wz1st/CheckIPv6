#!/bin/sh

cp -r CheckIPv6 /usr/local/bin/
chmod +x /usr/local/bin/CheckIPv6

cp -r check-server.service /etc/systemd/system/

systemctl daemon-reload
systemctl enable check-server.service
systemctl start check-server.service