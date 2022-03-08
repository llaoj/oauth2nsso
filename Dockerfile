FROM golang:1.17

WORKDIR /go/src/app
COPY . .
# 测试配置
COPY config.example.yaml  /etc/oauth2nsso/config.yaml

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["oauth2nsso", "--config=/etc/oauth2nsso/config.yaml"]