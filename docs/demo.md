## 单点登录演示

现在有两个应用`测试应用1`和`测试应用2`

### 登录

[测试应用1](http://oauth2nssodemo.p.rutron.net/authorize?client_id=test_client_1&response_type=code&scope=all&state=xyz&redirect_uri=http://localhost:9093/cb)

[测试应用2](http://oauth2nssodemo.p.rutron.net/authorize?client_id=test_client_2&response_type=code&scope=all&state=xyz&redirect_uri=http://localhost:9094/cb)

点击其中一个应用, **要求登录**, 输入`admin/admin`实完成登录, OAuth2NSSO会回掉返回`code`, 使用该`code`就可以调用 [接口1-2](../README.md#1-2-%E4%BD%BF%E7%94%A8code%E4%BA%A4%E6%8D%A2token) 获取`access_token`.

再次点击另外一个应用, **无需登录**, 直接回掉返回`code`. 同样, 使用`code`该应用可以调用 [接口1-2](../README.md#1-2-%E4%BD%BF%E7%94%A8code%E4%BA%A4%E6%8D%A2token) 获取`access_token`.

### 退出

在`测试应用1`或`测试应用2`中都可以调用退出逻辑

[点击退出](http://oauth2nssodemo.p.rutron.net/logout?redirect_uri=http%3A%2F%2Foauth2nssodemo.p.rutron.net%2Fauthorize%3Fclient_id%3Dtest_client_1%26response_type%3Dtoken%26scope%3Dall%26state%3Dxyz%26redirect_uri%3Dhttp%3A%2F%2Flocalhost%3A9093%2Fcb)

点击退出按钮, 会重新回到登录页面, 该服务退出后, 再次点击`测试应用1/测试应用2`会要求登录.