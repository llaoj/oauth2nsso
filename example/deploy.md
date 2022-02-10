## 示例程序部署命令

单机容器部署

```sh
docker run --restart always -d --name oauth2nssodemo -p 8083:9096 registry.cn-beijing.aliyuncs.com/llaoj/oauth2:0.2.0
```

## 请求示例

获取code请求连接:

http://oauth2nssodemo.p.rutron.net/authorize?client_id=test_client_1&response_type=token&scope=all&state=xyz&redirect_uri=http://localhost:9093/cb