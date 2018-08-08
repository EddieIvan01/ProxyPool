# ProxyPool

**Demo**

![](https://github.com/EddieIvan01/ProxyPool/blob/master/demo.png)

***

**流程**：

爬取免费代理，验证可用性存入本地sqlite3数据库，并在本地开启json API Server监听2333端口

***

**API**

+ /index：判断是否监听成功
+ /proxy：操作代理数据
  + ?act=get：获取数据库全部代理
  + ?act=reflush：刷新数据库
+ 状态码：
  + 200：正常
  + 201：数据库为空，调用reflush
  + 202：请求参数错误
+ 数据形式：

anonymous: 是否高匿

ssl：是否https

```
data{
    "code" : 200,
    "proxies" : [
        {
            "ip" : "118.190.95.43",
            "port" : "9001",
            "anonymous" : "1",
            "ssl" : "0"
        },
        {
            "ip" : "171.39.1.149",
            "port" : "8123",
            "anonymous" : "1",
            "ssl" : "1"
        }
        .
        .
        .
    ]
}
```

***

**架构**：

判断本地是否已存在数据库:

+ =>True: 返回数据库句柄，pass
+ =>False: 建立数据库，进行创建表等初始化工作，请求代理网站并筛选，将IP插入数据库，返回句柄

开启监听，提供API:

+ /proxy?get: 查询数据库，并返回数据库中全部代理
+ /proxy?reflush: 检查数据库已存在代理的可用性；请求代理网站获取IP，将可用IP插入数据库

------

**验证策略**：

请求[chinaz](http://ip.chinaz.com/getip.aspx)

+ http代理: 代理IP与响应IP相同 && 请求成功 && 请求响应 <= 5s

+ https代理: 请求成功 && 请求响应 <= 5s

***

**代理**：

- [西刺](http://www.xicidaili.com)
- [66ip](http://www.66ip.cn)
- [proxylist](https://list.proxylistplus.com)
