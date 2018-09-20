#!/bin/sh
GOOS=linux go build
tar -czf wechat-notify.tar.gz wechat-notify
curl -s fs.qianbao-inc.com/k8s/soft/uploadapi -F file=@wechat-notify.tar.gz -F truncate=yes
cksum ./wechat-notify