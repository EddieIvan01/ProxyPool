/*
*Author: Eddie_Ivan
*Blog: http://nemesisly.xyz
*Github: https://github.com/eddieivan01
*
*爬虫代理IP池
*抓取
*    http://www.xicidaili.com
*    http://www.66ip.cn
*    https://list.proxylistplus.com
*三家代理IP，存入本地Sqlite3数据库中，并在本地开启Http Server监听，提供json API服务
*
*程序架构：
*+ 判断本地是否已存在数据库
*  + =>True: 返回数据库句柄，pass
*  + =>False: 建立数据库，进行创建表等初始化工作，请求代理网站并筛选，将IP插入数据库，返回句柄
*+ 开启监听，提供API
*  + /proxy?act=get: 查询数据库，并返回数据库中全部代理
*  + /proxy?act=reflush: 检查数据库已存在代理的可用性；请求代理网站获取IP，将可用IP插入数据库
*
*验证策略：
*请求http://ip.chinaz.com/getip.aspx
*http代理: 代理IP与响应IP相同 && 请求成功 && 请求响应 <= 5s
*https代理: 请求成功 && 请求响应 <= 5s
*
*API状态码：
* 200：正常
* 201：数据库为空
* 202：请求参数错误
*
 */
package main

import "fmt"

func main() {

	fmt.Println(`
       ___                            ___              __
      / _ \  ____ ___  __ __  __ __  / _ \ ___  ___   / /
     / ___/ / __// _ \ \ \ / / // / / ___// _ \/ _ \ / / 
    /_/    /_/   \___//_\_\  \_, / /_/    \___/\___//_/  
                            /___/                        
    `)

	initHTTP()
}
