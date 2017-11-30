# zb

[![Build Status](https://travis-ci.org/shchnmz/zb.svg?branch=master)](https://travis-ci.org/shchnmz/zb)
[![Go Report Card](https://goreportcard.com/badge/github.com/shchnmz/zb)](https://goreportcard.com/report/github.com/shchnmz/zb)
[![GoDoc](https://godoc.org/github.com/shchnmz/zb?status.svg)](https://godoc.org/github.com/shchnmz/zb)

zb是一个基于[明日系统导出至redis中的学校数据](https://github.com/shchnmz/ming)，提供转班相应操作的一个[Golang](https://golang.org)包。

#### 服务器实现
* [api](./api)是API服务器。
* [server](./server)是应用服务器,提供了用户转班申请界面以及简单的admin管理页面。

#### 生产环境搭建
* Redis安装
  [Install and Configure Redis on CentOS 7](https://github.com/northbright/Notes/blob/master/Redis/Install/Install_and_Config_Redis_on_CentOS.md)
* Nginx反向代理设置

        // sudo vi /etc/nginx/nginx.conf

        server {
            listen 80;
            server_name localhost;

            location / {
                proxy_pass http://localhost:8080;
            }

            location /api {
                proxy_pass http://localhost:8000;
            }
        }

* API服务器，应用服务器作为系统服务自动启动
  [Configure Binary as systemd Service on CentOS 7](https://github.com/northbright/Notes/blob/master/Linux/CentOS/service/config-binary-as-systemd-service-on-centos-7/config-binary-as-systemd-service-on-centos-7.md)

#### Documentation
* [API References](https://godoc.org/github.com/shchnmz/zb)

#### License
* [MIT License](LICENSE)
