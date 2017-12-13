# server

server是提供转班服务以及管理的HTTP服务器。它是使用[Golang](https://golang.org)写的。

#### 如何使用
1. 确认已经运行过[ming800-to-redis](https://github.com/shchnmz/ming/tree/master/tools/ming800-to-redis)将明日系统的数据导入到redis中。

2. 在和可执行文件相同的文件夹下，创建一个`config.json`的配置文件：

        {
            "server_addr": ":8080",
            "redis_server": "localhost:6379",
            "redis_password": "",
            "admin_account": "admin",
            "admin_password": "admin",
            "closed_notices":[
               "在线转班申请已截止，请等待学校电话通知处理结果。",
               "联系电话:xxxx。"
            ]
        }

* `"server_addr"`是应用服务器的运行地址。
* `"redis_server"`,`"redis_password"`是同步的redis的地址和密码。
* `"account_account`,`"account_password"`是HTTP Basic Auth的用户名和密码，用来管理／查看后台。
* `"closed_notices"`是关闭转班系统后的提示信息，可以有多个(多行)提示信息。

3. 运行`server`

        ./server

4. 路由
* 学生访问
  * `/`

    学生登录入口，登记转班信息。

* 管理员访问

  需要config.json中的`"admin_account"`和`"admin_password"`来进行HTTP Basic Auth进行查看。
  * `/admin`

    管理员查看所有转班申请记录。

  * `/statistics`

    管理员查看转班信息统计。
    包括如下统计信息：
      * 按校区

        从某某校区转入某某校区的申请转班学生人数
      * 按教师转班率

        转班率 = 某一个老师所教授的学生中，转班申请学生人数 / 某一个老师所教授所有学生人数
      * 按教师

        列出某一个老师所教授的学生中，转班申请学生人数

  * `/enable`

    允许学生进行转班申请。

  * `/disable`

    关闭转班系统，不允许学生进行转班申请。

#### 配置禁止转班的信息
* [从配置文件中加载转班禁止信息到redis中](../tools/load-blacklist-from-file)
* [从redis数据库中清除转班禁止信息](../tools/clear-blacklist)
