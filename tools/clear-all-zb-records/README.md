# clear-all-zb-records

clear-all-zb-records是一个从redis数据库中清除所有转班记录的程序。它是使用[Golang](https://golang.org)写的。

#### 如何使用
1. 确认已经运行过[ming800-to-redis](https://github.com/shchnmz/ming/tree/master/tools/ming800-to-redis)将明日系统的数据导入到redis中。

2. 在和可执行文件相同的文件夹下，创建一个`config.json`的配置文件：

        {
            "redis_server": "localhost:6379",
            "redis_password": ""
        }

* `"redis_server"`,`"redis_password"`是同步的redis的地址和密码。

3. 运行`clear-all-zb-records`

        ./clear-all-zb-records
