FROM golang:1.17

WORKDIR /go/src/app
COPY . .
# 测试配置
COPY config.example.yaml  /etc/oauth2/config.yaml

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["oauth2", "--config=/etc/oauth2/config.yaml"]