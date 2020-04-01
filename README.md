# oauth2
oauth2 server based on go-oauth2
实现了auth2的四种工作流程
1. authorization_code
2.
3.
4.
5. 验证access_token
---

## 1. Flow: authorization_code

### 1.1. 获取授权code

**方法**  
GET

**请求示例**  
`http://localhost:9096/authorize?client_id=222222&response_type=code&scope=all&state=xyz&redirect_uri=http://localhost:9094/cb`


**参数说明**  

|参数|类型|说明|
|-|-|-|
|client_id|string|在oauth2 server 注册的client_id|
|response_type|string|固定值`code`|
|scope|string|权限范围,`str1,str2,str3`|
|state|string|验证请求的标志字段,后续第二部需要先验证该字段是否是第一步设置的值|
|redirect_uri|string|发放`code`用的回调uri,回调时会在uri后面跟上`?code=**&state=###`|

**返回示例**  
`302 http://localhost:9094/cb?code=XUNKO4OPPROWAPFKEWNZWA&state=xyz`

### 1.2. 使用`code`交换`token`

**Method**  
POST

**Url**  
`http://localhost:9096/token`

**Authorization**
- basic auth
- username: `client_id`
- password: `client_secret`

**Header**  
`Content-Type: application/x-www-form-urlencoded`

**Body参数说明**  

|参数|类型|说明|
|-|-|-|
|grant_type|string|固定值`authorization_code`|
|code|string|第一步发放的code|
|redirect_uri|string|第一步填写的redirect_uri|

**Response返回示例**  

```
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIyMjIyMjIiLCJleHAiOjE1ODU3MTU1NTksInN1YiI6InRlc3QifQ.ZMgIDQMW7FGxbF1V8zWOmEkmB7aLH1suGYjhDdrT7aCYMEudWUoiCkWHSvBmJahGm0RDXa3IyDoGFxeMfzlDNQ",
    "expires_in": 7200,
    "refresh_token": "JG7_WGLWXUOW2KV2VLJKSG",
    "scope": "all",
    "token_type": "Bearer"
}
```

## 5. 验证token

**接口说明**  
这个接口是资源服务端使用的, 用来验证token, **如果返回正确, 资源服务端还要验证下scope**

**请求方法**  
GET

**Url**  
`http://localhost:9096/test`

**Authorization**
- Bearer Token
- Token: `access_token`

**返回示例**  

```
Status Code: 200
Response Body
{
  "client_id": "222222",
  "expires_in": 7191,
  "scope": "all",
  "user_id": "test"
}
```

```
Status Code: 400
Response Body
   invalid access token
```