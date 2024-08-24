#!/bin/bash
echo "--*-- 卸载开始..."

systemctl disable skpy
echo "  |--> 删除服务"
systemctl stop skpy
rm -f /etc/systemd/system/skpy.service

echo "  |--> 删除文件"
rm -rf /opt/my-apps/skpy
rm -rf /var/log/skpy

echo "  |--> 服务刷新"
systemctl daemon-reload

echo "--*-- 卸载完成 ^_^"
