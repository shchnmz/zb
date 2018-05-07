# gen-blacklist-json-from-txt

gen-blacklist-json-from-txt是一个将文本TXT中的转班禁止配置整合生成JSON格式的配置文件。它是使用[Golang](https://golang.org)写的。

#### 转班禁止配置TXT文件以及格式
* 在可执行文件相同文件夹下，创建一个`blacklist`文件夹
* 其中可以存放6个文件
  * `from_campuses.txt` - 禁止转出的校区
  * `from_periods.txt` - 禁止转出的时段 
  * `from_classes.txt` - 禁止转出的班级
  * `to_campuses.txt` - 禁止转入的校区
  * `to_periods.txt` - 禁止转入的时段
  * `to_classes.txt` - 禁止转入的班级
* 每一行为1条记录

        // e.g.

        // from_campuses.txt
        校区D

        // to_campuses.txt
        校区D

        // from_classes.txt
        校区B:二年级:二年级3班
        校区A:三年级:三年级1班

        // to_periods.txt
        校区A:五年级:星期五16:25-17:55
        校区B:五年级:星期六16:25-17:55

#### 最后生成的转班禁止配置JSON文件以及格式
* 在和可执行文件相同文件夹下，运行后会生成`blacklist.json`的配置文件:

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

#### 如何使用
1. 在可执行文件相同文件夹下，创建一个`blacklist`文件夹
2. 运行[list-all-classes](https://github.com/shchnmz/ming/tree/master/tools/list-all-classes)和[list-all-periods](https://github.com/shchnmz/ming/tree/master/tools/list-all-periods)生成所有班级和时段列表，在此基础上修改生成禁止转班的TXT文件 
3. 运行`gen-blacklist-json-from-txt`

        ./gen-blacklist-json-from-txt
4. 检查当前目录下的`blacklist.json`
