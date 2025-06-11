#!/bin/sh

cp -r check_ipv6.sh /etc/config
chmod +x /etc/config/check_ipv6.sh

DEFAULT_IP=""
IP="${1:-$DEFAULT_IP}"

CRON_JOB="*/1 * * * * /etc/config/check_ipv6.sh $IP >/dev/null 2>&1"
(crontab -l 2>/dev/null | grep -F -q "$CRON_JOB") || (
    (crontab -l 2>/dev/null; echo "$CRON_JOB") | crontab -
    echo "已将任务写入 crontab。"
)
/etc/init.d/cron restart