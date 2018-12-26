FROM alpine
WORKDIR /app
COPY wechat-notify /app/wechat-notify
#ENTRYPOINT /app/wechat-notify
CMD /app/wechat-notify -h
EXPOSE 8001