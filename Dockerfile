FROM harbor.haodai.net/base/alpine:3.7cgo
WORKDIR /app

MAINTAINER wenzhenglin(http://g.haodai.net/wenzhenglin/wechat-notify.git)

COPY wechat-notify /app
CMD /app/wechat-notify -h
ENTRYPOINT /app/wechat-notify

EXPOSE 8001