CGO_ENABLED=0 go build
docker build -t wechat-notify .
docker tag wechat-notify harbor.haodai.net/ops/wechat-notify:v1
docker push harbor.haodai.net/ops/wechat-notify:v1