#! /bin/bash
echo "--*-- 安装开始..."
mkdir -p /opt/my-apps
cp -rf files/skpy /opt/my-apps

echo "  |--> 安装服务"
cp -f files/skpy.service /usr/lib/systemd/system


echo "  |--> 服务刷新"
systemctl daemon-reload

echo "  |--> 启动服务"
systemctl start skpy

echo "  |--> 服务开机启动"
systemctl enable skpy

systemctl status skpy

echo "--*-- 安装完成 ^_^"




