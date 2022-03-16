# OAuth2&SSO

## 项目介绍

[llaoj/oauth2nsso](https://github.com/llaoj/oauth2nsso) 项目是基于 go-oauth2 打造的**独立**的 OAuth2.0 和 SSO 服务，提供了开箱即用的 OAuth2.0服务和单点登录SSO服务。开源一年多，获得了社区很多用户的关注，该项目多公司线上在用，其中包含上市公司。轻又好用，稳的一P。

感谢:
![sponsors](https://raw.githubusercontent.com/llaoj/oauth2nsso/master/docs/sponsors.png)

## B站视频讲解

 [教你构建OAuth2.0和SSO单点登录服务(基于go-oauth2)](https://www.bilibili.com/video/BV1UA411v73P)

## 单点登录(SSO)示例

[单点登录(SSO)示例](http://rutron.net/docs/oauth2nsso/demo/)

## 动图演示

授权码(authorization_code)流程 & 单点登录(SSO)

![authorization_code_n_sso](https://raw.githubusercontent.com/llaoj/oauth2nsso/master/docs/demo-pic/authorization_code_n_sso.gif)

## 主要功能
https://raw.githubusercontent.com/llaoj/oauth2nsso/master/docs/demo-pic/authorization_code_n_sso.gif
**实现了oauth2的四种工作流程**

1. authorization_code
2. implicit
3. password
4. client credentials

**扩展功能**

5. 资源端用的验证 access_token 接口 `/verify`
6. 刷新 token 接口 `/refresh`
7. 专门为 SSO 开发的客户端登出接口 `/logout`

详情见[API说明](http://rutron.net/docs/oauth2nsso/apis/)


# 配置

该项目的配置修改都是在配置文件中完成的，配置文件在启动应用的时候通过`--config=`标签进行配置.

配置文件介绍如下：

```yaml
# session 相关配置
session:
  name: session_id
  secret_key: "kkoiybh1ah6rbh0"
  # 过期时间
  # 单位秒
  # 默认20分钟
  max_age: 1200

# 用户登录验证方式
# 支持: db ldap
auth_mode: ldap

# 数据库相关配置
# 这里可以添加多个连接支持
# 默认是 default 连接
db:
  default:
    type: mysql
    host: string
    port: 3306
    user: 123
    password: abc
    dbname: oauth2nsso

ldap:
  # 服务地址
  # 支持 ldap ldaps
  url: ldap://ldap.forumsys.com
  # url: ldaps://ldap.rutron.net

  # 查询使用的DN
  search_dn: cn=read-only-admin,dc=example,dc=com
  # 查询使用DN的密码
  search_password: password
  
  # 基础DN
  # 以此为基础开始查找用户
  base_dn: dc=example,dc=com
  # 查询用户的Filter
  # 比如: 
  #   (&(uid=%s)) 
  #   或 (&(objectClass=organizationalPerson)(uid=%s))
  #   其中, (uid=%s) 表示使用 uid 属性检索用户, 
  #   %s 为用户名, 这一段必须要有, 可以替换 uid 以使用其他属性检索用户名
  filter: (&(uid=%s))

# 可选
# redis 相关配置
# 可以提供:
# - 统一回话存储
# - oauth2 client 存储
redis:
  default:
    addr: 127.0.0.1:6379
    password: 
    db: 0

# oauth2 相关配置
oauth2:
  # access_token 过期时间
  # 单位小时
  # 默认2小时
  access_token_exp: 2
  # 签名 jwt access_token 时所用 key
  jwt_signed_key: "k2bjI75JJHolp0i"
  
  # oauth2 客户端配置
  # 数组类型
  # 可配置多客户端
  client:

      # 客户端id 必须全局唯一
    - id: test_client_1
      # 客户端 secret
      secret: test_secret_1
      # 应用名 在页面上必要时进行显示
      name: 测试应用1
      # 客户端 domain
      # !!注意 http/https 不要写错!!
      domain: http://localhost:9093
      # 权限范围
      # 数组类型
      # 可以配置多个权限 
      # 颁发的 access_token 中会包含该值 资源方可以对该值进行验证
      scope:
          # 权限范围 id 唯一
        - id: all
          # 权限范围名称
          # 会在页面（登录页面）进行展示
          title: "用户账号、手机、权限、角色等信息"

    - id: test_client_2
      secret: test_secret_2
      name: 测试应用2 
      domain: http://localhost:9094
      scope:
        - id: all
          title: 用户账号, 手机, 权限, 角色等信息

```


# API列表

## 1 authorization_code

### 1-1 获取授权code

**请求方式**

`GET` `/authorize`

**参数说明**  

|参数|类型|说明|
|-|-|-|
|client_id|string|在oauth2 server注册的client_id,见配置文件[oauth2.client.id](http://rutron.net/docs/oauth2nsso/configuration/)|
|response_type|string|固定值:`code`|
|scope|string|权限范围,如:`str1,str2,str3`,str为配置文件中[oauth2.client.scope.id](http://rutron.net/docs/oauth2nsso/configuration/)的值 |
|state|string|表示客户端的当前状态,可以指定任意值,认证服务器会原封不动地返回这个值|
|redirect_uri|string|回调uri,会在后面添加query参数`?code=xxx&state=xxx`,发放的code就在其中|

**请求示例**

```sh
# 浏览器请求
http://localhost:9096/authorize?client_id=test_client_1&response_type=code&scope=all&state=xyz&redirect_uri=http://localhost:9093/cb

# 302跳转,返回code
http://localhost:9093/cb?code=XUNKO4OPPROWAPFKEWNZWA&state=xyz
```

### 1-2 使用`code`交换`token`

**请求方式**

`POST` `/token`

**请求头 Authorization**

- basic auth
- username: `client_id`
- password: `client_secret`

**Header**  
`Content-Type: application/x-www-form-urlencoded`

**Body参数说明**  

|参数|类型|说明|
|-|-|-|
|grant_type|string|固定值`authorization_code`|
|code|string| 1-1 发放的code|
|redirect_uri|string| 1-1 填写的redirect_uri|

**Response返回示例**  

```json
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIyMjIyMjIiLCJleHAiOjE1ODU3MTU1NTksInN1YiI6InRlc3QifQ.ZMgIDQMW7FGxbF1V8zWOmEkmB7aLH1suGYjhDdrT7aCYMEudWUoiCkWHSvBmJahGm0RDXa3IyDoGFxeMfzlDNQ",
    "expires_in": 7200,
    "refresh_token": "JG7_WGLWXUOW2KV2VLJKSG",
    "scope": "all",
    "token_type": "Bearer"
}
```

## 2 implicit

资源请求方(client方)使用, 
多用于没有后端的应用, 
用户授权登录之后, 会直接向前端发送令牌(`access_token`)

**请求方式**

`GET` `/authorize`

**参数说明**  

|参数|类型|说明|
|-|-|-|
|client_id|string|在 oauth2 server 注册的client_id|
|response_type|string|固定值`token`|
|scope|string|权限范围,同1-1中说明|
|state|string|验证请求的标志字段|
|redirect_uri|string|发放`code`用的回调uri,回调时会在uri后面跟上`?code=**&state=###`|

**请求示例**

```sh
http://localhost:9096/authorize?client_id=test_client_1&response_type=token&scope=all&state=xyz&redirect_uri=http://localhost:9093/cb
```

**返回示例**

```sh
# 302 跳转,返回 access_token
http://localhost:9093/cb#access_token=eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0X2NsaWVudF8xIiwiZXhwIjoxNTkxNDI3OTMwLCJzdWIiOiJhZG1pbiJ9.RBYns9UnNYDHINSBzvHWHRzuKCpzKmsxUnKt30lntmGvXmVDoByZtlB0RHAVB59PHBlJNO_YUBZzC2odwCa8Tg
&expires_in=3600&scope=all&state=xyz&token_type=Bearer
```

**注意**

1. 这里会返回请求时设置的`state`, 请在进行下一步之前验证它, 防止请求被劫持或者篡改
2. 这种方式把令牌直接传给前端，是很不安全的。因此，只能用于一些安全要求不高的场景，并且令牌的有效期必须非常短，通常就是会话期间（session）有效，浏览器关掉，令牌就失效了


## 3 password

资源请求方(client方)使用
如果充分信任接入应用(client), 用户就可以直接把用户名密码给接入应用.
接入应用使用用户账号密码申请令牌.

**请求方式**

`POST` `/token`

**请求头 Authorization**

- basic auth
- username: `client_id`
- password: `client_secret`

**Header**  
`Content-Type: application/x-www-form-urlencoded`

**Body参数说明**  

|参数|类型|说明|
|-|-|-|
|grant_type|string|固定值`password`|
|username|string|用户名|
|password|string|用户密码|
|scope|string|权限范围,同1-1中说明|

**返回示例**  

```json
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0X2NsaWVudF8xIiwiZXhwIjoxNTkxNDMyNzA3LCJzdWIiOiJhZG1pbiJ9.ECfUkCMUZE8I6GH3XTDcJnQgDryiRyyBhEHBW-dCxzFWaR-mvU5dsx3XV2bx-LWZzPJTBAQ3rB5QOb4BHjnBXw",
    "expires_in": 7200,
    "refresh_token": "AH-B00RKXTME9WXDPSBSTG",
    "scope": "all",
    "token_type": "Bearer"
}
```

## 4 client_credentials

资源请求方(client方)使用
使用在oauth2服务器注册的client_id 和 client_secret 获取 access_token,
发出 API 请求时，它应将access_token作为 Bearer 令牌传递到 Authorization 请求头中。

**请求方式**

`POST` `/token`

**请求头 Authorization**

- basic auth
- username: `client_id`
- password: `client_secret`

**Header**  

`Content-Type: application/x-www-form-urlencoded`

**Body参数说明**  

|参数|类型|说明|
|-|-|-|
|grant_type|string|固定值`client_credentials`|
|scope|string|权限范围,同1-1中说明|

**返回示例**  

```json
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJlbWJlZGVkLWg1LWFwaSIsImV4cCI6MTU4OTk3NzQyNn0.Pu93fy0-gyiFqExBkCFAKTVJ1on_RpOSexzkHqczA6n6kB2_mOHbTMOyGK_Di7bHxZ3JqpZeyDoKQBtUe_T7jw",
    "expires_in": 7200,
    "token_type": "Bearer"
}
```

## 5 验证token

**接口说明**

这个接口是资源端使用的, 
用来验证 `access_token` `scope` 和 `domain` .

**请求方式**

`GET`  `/verify`

**请求头 Authorization**

- Bearer Token
- Token: `access_token`

**返回示例**  

正确 Status Code: 200

Response Body:

```json
{
  "client_id": "test_client_1",
  "domain": "http://127.0.0.1:9093",
  "expires_in": 7188,
  "scope": "all",
  "user_id": ""
}
```
> **注意:** 接口还会一起返回权限范围`scope` 和 client 的注册 domain, 这里推荐验证下, 请求方的身份和权限

错误 Status Code: 400

Response Body: `invalid access token`

## 6 刷新token

刷新access_token, 使用refresh_token换取access_token

**请求方式**

`POST` `/token`

**请求头 Authorization**

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

```json
{
    "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIyMjIyMjIiLCJleHAiOjE1ODU4MTc2MTMsInN1YiI6IjEifQ.yNpQIbklhtsSr5KEkJMAR4I30c85OEriYwAOpL_ukRBJ1qsSziT05HFN-kxVN1-qM18TzVEf8beCvugyhpgpsg",
    "expires_in": 7200,
    "refresh_token": "2AH_LQHPUYK8XML4LKMQKG",
    "scope": "all",
    "token_type": "Bearer"
}
```

## 7 logout

专门为SSO开发, 
主要是销毁浏览器的会话, 退出登录状态, 跳转到指定链接(redirect_uri)

**请求方式**

`GET` `/logout?redirect_uri=xxx`

**参数说明**  

|参数|类型|说明|
|-|-|-|
|redirect_uri|string|退出登录后跳转到的地址,建议使用1-1所生成的地址, 需要urlencode|

**请求示例**  

```sh
http://localhost:9096/logout?redirect_uri=http%3a%2f%2flocalhost%3a9096%2fauthorize%3fclient_id%3dtest_client_1%26response_type%3dcode%26scope%3dall%26state%3dxyz%26redirect_uri%3dhttp%3a%2f%2flocalhost%3a9093%2fcb
```


# 部署

## 修改配置和完善代码

克隆到代码之后，首先需要进行配置文件的修改和部分代码逻辑的编写：

```sh
# 克隆源码
git clone git@github.com:llaoj/oauth2nsso.git
cd oauth2nsso

# 根据实际情况修改配置
cp config.example.yaml /etc/oauth2nsso/config.yaml
vi /etc/oauth2nsso/config.yaml
...

# 如果使用 LDAP方式 验证用户, 直接修改配置文件即可
# OR
# 如果使用 数据库方式 验证用户, 需要修改源码
# 主要修改登录部分逻辑:
# 文件: model/user.go:21
# 方法: Authentication()
...
```

## 使用docker部署

**[推荐]** 容器化部署比较方便进行大规模部署，是当下的趋势。需要本地有 docker 环境。

```sh
# 构建镜像
docker build -t <image:tag> .

# 运行
docker run --rm --name=oauth2nsso --restart=always -d \
-p 9096:9096 \
-v <path to config.yaml>:/etc/oauth2nsso/config.yaml \
<image:tag>
```

## 基于源码部署

```sh
# 在仓库根目录
# 编译
go build -mod=vendor

# 运行
./oauth2nsso -config=/etc/oauth2nsso/config.yaml
```


# 版本说明

## v0.2.0

该项目发布以来收到了很多朋友的关注，很多公司都将它应用到了一些比较重要的项目中。同时，也对该项目提出了很多要求。综合这些，开发了这个版本。同时希望朋友们互相交流，多提意见。

这个版本主要有下面几个改动：

1. 由于 go-oauth2.v3 版本安全性原因，将该包升级到 v4
2. 丰富了可配置的项目
3. 增加了容器化部署的脚本和相关文档
4. 多了一些细节的优化
5. 增加了错误页面
6. 用户验证增加了LDAP支持
