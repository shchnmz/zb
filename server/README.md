# server

server是提供转班服务以及管理的HTTP服务器。它是使用[Golang](https://golang.org)写的。

#### 如何使用
1. 确认已经运行过[ming800-to-redis](../ming800-to-redis)将明日系统的数据导入到redis中。

2. 在和可执行文件相同的文件夹下，创建一个`config.json`的配置文件：

        {
            "server_addr": ":8080",
            "redis_server": "localhost:6379",
            "redis_password": "",
            "admin_account": "admin",
            "admin_password": "admin",
        }

* `"server_addr"`是应用服务器的运行地址。
* `"redis_server"`,`"redis_password"`是同步的redis的地址和密码。
* `"account_account`,`"account_password"`是HTTP Basic Auth的用户名和密码，用来管理／查看后台。

3. 运行`server`

        ./server

