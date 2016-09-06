# tomato

Parse-compatible API server module for Golang/Beego

## 开始
###### 安装
```bash
    go get github.com/lfq7413/tomato
```
###### 创建文件 hello.go
```go
package main

import "github.com/lfq7413/tomato"

func main() {
    tomato.Run()
}
```
###### 创建配置文件 /conf/app.conf
```ini
appname = hello
httpport = 8080
runmode = dev
autorender = false
copyrequestbody = true
EnableDocs = true

serverurl = http://127.0.0.1:8080/v1/
databaseuri = 127.0.0.1:27017/test
```
###### 运行
```bash
    go run hello.go
```
###### 完成
打开浏览器，访问 `http://127.0.0.1:8080`

## 功能

## 注意事项

* 修改 tomato 工程的路由信息之后，需要运行 refresh.sh ，更新注解路由
* 使用 tomato 的应用，可在应用同目录下添加配置文件：/conf/app.conf，文件格式与 beego 完全一致

## 开发日志

[开发日志.md](/开发日志.md)

## LICENSE

[MIT](/LICENSE)
