## 示例程序部署命令

单机容器部署

```sh
docker run --restart always -d --name oauth2nssodemo -p 8083:9096 registry.cn-beijing.aliyuncs.com/llaoj/oauth2:0.2.0
```

## 前置LB(NGINX)配置

oauth2nsso 服务前置负载均衡配置文件
根据自己情况修改

```sh
upstream backend_oauth2nssodemo_p_rutron_net {
    server 172.19.168.60:8083 weight=1 max_fails=3 fail_timeout=30s;
}

server {
    listen 80;
    server_name oauth2nssodemo.p.rutron.net;

    charset utf-8;
    client_max_body_size 120m;

    proxy_connect_timeout 180;
    proxy_send_timeout 180;
    proxy_read_timeout 180;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;

    location / {
        proxy_pass http://backend_oauth2nssodemo_p_rutron_net;
    }
}
```

## 请求示例

获取code请求连接:

http://oauth2nssodemo.p.rutron.net/authorize?client_id=test_client_1&response_type=token&scope=all&state=xyz&redirect_uri=http://localhost:9093/cb