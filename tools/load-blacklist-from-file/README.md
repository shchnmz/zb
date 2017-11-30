# load-blacklist-from-file

load-blacklist-from-file是一个从配置文件中加载转班禁止信息到redis中的程序。它是使用[Golang](https://golang.org)写的。

#### 转班禁止配置文件以及格式
* 在和可执行文件相同文件夹下，创建一个`blacklist.json`的配置文件:

        {
          "blacklist": {
            "from_campuses":["校区D"],
            "from_periods":[],
            "from_classes":[
                "校区B:二年级:二年级3班",
                "校区A:三年级:三年级1班"
            ],
            "to_campuses":["校区D"],
            "to_periods":[
                "校区A:五年级:星期五16:25-17:55",
                "校区B:五年级:星期六16:25-17:55"
            ],
            "to_classes":[]
        }

  * `"from_campuses"`,`"to_campuses"`是禁止转出／转入的校区列表。
  * `"from_periods"`,`"to_periods"`是禁止转出／转入的时间段列表。

     时间段格式为:`$校区:$课程:$时间段`。 e.g. "校区A:五年级:星期五16:25-17:55"。

  * `"from_classes"`, `"to_classes"`是禁止转出／转入的班级列表。

     班级格式为:`$校区:$课程:$班级`。 e.g. "校区B:二年级:二年级3班"。 

#### 如何使用
1. 确认已经运行过[ming800-to-redis](../ming800-to-redis)将明日系统的数据导入到redis中。

2. 在和可执行文件相同的文件夹下，创建一个`config.json`的配置文件：

        {
            "redis_server": "localhost:6379",
            "redis_password": ""
        }

* `"redis_server"`,`"redis_password"`是同步的redis的地址和密码。

3. 运行`load-blacklist-from-file`

        ./load-blacklist-from-file
