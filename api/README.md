# api

api是提供转班服务API的HTTP服务器。它是使用[Golang](https://golang.org)写的。

## 部署
1. 确认已经运行过[ming800-to-redis](https://github.com/shchnmz/ming/tree/master/tools/ming800-to-redis)将明日系统的数据导入到redis中。

2. 在和可执行文件相同的文件夹下，创建一个`config.json`的配置文件：

        {
            "server_addr": ":8000",
            "redis_server": "localhost:6379",
            "redis_password": ""
        }

* `"server_addr"`是API服务器的运行地址。
* `"redis_server"`,`"redis_password"`是同步的redis的地址和密码。

3. 运行`api`

        ./api


## API文档 

#### 根据手机号码获取绑定的学生
返回手机号码绑定的学生姓名列表。
一个手机号码可以绑定多个学生。

    GET /api/get-names-by-phone-num/:phone_num

    status: 200 OK

    {
        "success":true,
        "err_msg":"",
        "names":[
            "学生A",
            "学生B"
        ]
    }

#### 根据学生信息（姓名，手机号码的组合）获取学生就读班级
返回学生（姓名，手机号码）所就读的班级信息列表。
一个学生可以就读多个班级。
班级信息的格式：`校区:课程:班级`。e.g. `校区A:一年级绘画课程:一年级（1）班`

    GET /api/get-classes-by-name-and-phone-num/:name/:phone_num

    status: 200 OK

    {
        "success":true，
        "err_msg":"",
        "classes":[
            "校区A:二年级:二年级6",
            "校区A:2018暑假零基础班:18暑新二三年级3"
        ]
    }

#### 根据班级信息获取授课教师
返回班级的授课教师姓名列表。
一个班级可以有多个授课教师。
班级信息的格式：`校区:课程:班级`。e.g. `校区A:一年级绘画课程:一年级（1）班`

    GET /api/get-teachers-by-class/:class

    status: 200 OK

    {
        "success":true,
        "err_msg":"",
        "teachers":[
            "教师A",
            "教师B",
        ]
    }

