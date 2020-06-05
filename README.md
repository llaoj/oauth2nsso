# gooauth
1. oauth2 server: based on go-oauth2
2. sso: based on the oauth2 service

# 目录

**实现了auth2的四种工作流程**

1. [authorization_code](##1.Flow:authorization_code)
2. implicit
3. password
4. client credentials

**其他**

5. 验证access_token (资源端)
6. 刷新token
7. 专门为SSO开发的logout

**配置**

1. implicit 和 client credentials 模式是不会生成refresh token的, 刷新token时会删除原有的token重新发布新的token.
2. 每一种模式的配置详情如下:

```
var (
  DefaultCodeExp               = time.Minute * 10
  DefaultAuthorizeCodeTokenCfg = &Config{AccessTokenExp: time.Hour * 2, RefreshTokenExp: time.Hour * 24 * 3, IsGenerateRefresh: true}
  DefaultImplicitTokenCfg      = &Config{AccessTokenExp: time.Hour * 1}
  DefaultPasswordTokenCfg      = &Config{AccessTokenExp: time.Hour * 2, RefreshTokenExp: time.Hour * 24 * 7, IsGenerateRefresh: true}
  DefaultClientTokenCfg        = &Config{AccessTokenExp: time.Hour * 2}
  DefaultRefreshTokenCfg       = &RefreshingConfig{IsGenerateRefresh: true, IsRemoveAccess: true, IsRemoveRefreshing: true}
)
```

---

## 1.Flow:authorization_code

### 1.1. 获取授权code

**方法**

GET

**Url**

`/authorize`

**请求示例**

```
http://localhost:9096/authorize?client_id=test_client_1&response_type=code&scope=all&state=xyz&redirect_uri=https://localhost:9093/cb
```


**参数说明**  

|参数|类型|说明|
|-|-|-|
|client_id|string|在oauth2 server 注册的client_id|
|response_type|string|固定值`code`|
|scope|string|权限范围,`str1,str2,str3`, 如果没有特殊说明,填`all` |
|state|string|验证请求的标志字段|
|redirect_uri|string|发放`code`用的回调uri,回调时会在uri后面跟上`?code=**&state=###`|

**返回示例**

`302 http://localhost:9093/cb?code=XUNKO4OPPROWAPFKEWNZWA&state=xyz`

**注意**

这里会返回请求时设置的`state`, 请在进行下一步之前验证它

### 1.2. 使用`code`交换`token`

**Method**

POST

**Url**

`/token`

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

## 2.Flow: implicit

## 3.Flow: password

## 4.Flow: client credentials

使用在oauth2服务器注册的client_id 和 client_secret 获取 access_token,
发出 API 请求时，它应将access_token作为 Bearer 令牌传递到 Authorization 请求头中。

**请求方法**

POST

**Url**

`/token`

**Authorization**

- basic auth
- username: `client_id`
- password: `client_secret`

**Header**  

`Content-Type: application/x-www-form-urlencoded`

**Body参数说明**  

|参数|类型|说明|
|-|-|-|
|grant_type|string|固定值`client_credentials`|
|scope|string|权限范围,`str1,str2,str3`, 如果没有特殊说明,填`all` |

**返回示例**  

```
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJlbWJlZGVkLWg1LWFwaSIsImV4cCI6MTU4OTk3NzQyNn0.Pu93fy0-gyiFqExBkCFAKTVJ1on_RpOSexzkHqczA6n6kB2_mOHbTMOyGK_Di7bHxZ3JqpZeyDoKQBtUe_T7jw",
    "expires_in": 7200,
    "token_type": "Bearer"
}
```


## 5. 验证token

**接口说明**

这个接口是资源端使用的, 用来验证 `access_token` 和 `scope` .

**请求方法**

GET

**Url**

`/test`

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

注意, 如果token正确, 还会一起返回权限范围`scope`, 这里需要验证下, 请求方是否拥有该权限.

```
Status Code: 400
Response Body
   invalid access token
```

## 6. 刷新token

刷新access_token, 使用refresh_token换取access_token

**Method**

POST

**Url**

`/token`

**Authorization**

- basic auth
- username: `client_id`
- password: `client_secret`

**Header**

`Content-Type: application/x-www-form-urlencoded`

**Body参数说明**

|参数|类型|说明|
|-|-|-|
|grant_type|string|固定值`refresh_token`|
|refresh_token|string|之前获取的refresh_token|

**返回示例**

```
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIyMjIyMjIiLCJleHAiOjE1ODU4MTc2MTMsInN1YiI6IjEifQ.yNpQIbklhtsSr5KEkJMAR4I30c85OEriYwAOpL_ukRBJ1qsSziT05HFN-kxVN1-qM18TzVEf8beCvugyhpgpsg",
    "expires_in": 7200,
    "refresh_token": "2AH_LQHPUYK8XML4LKMQKG",
    "scope": "all",
    "token_type": "Bearer"
}
```


## 7 logout

专门为SSO开发
主要是销毁浏览器的会话, 退出登录状态, 跳转到指定链接(redirect_uri)

**Method**

GET

**Url**

`/logout?redirect_uri=xxx`

**参数说明**  

|参数|类型|说明|
|-|-|-|
|redirect_uri|string|退出登录后跳转到的地址,建议使用1.1所生成的地址, 需要urlencode|

**请求示例**  

```
http://localhost:9096/logout?redirect_uri=http%3a%2f%2flocalhost%3a9096%2fauthorize%3fclient_id%3dtest_client_1%26response_type%3dcode%26scope%3dall%26state%3dxyz%26redirect_uri%3dhttp%3a%2f%2flocalhost%3a9093%2fcb
```