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
databaseuri = 192.168.99.100:27017/test
AppID = test
MasterKey = test
ClientKey = test
AllowClientClassCreation = true
```
###### 运行
```bash
    go run hello.go
```
###### 创建对象
```bash
    curl -X POST \
    -H "X-Parse-Application-Id: test" \
    -H "X-Parse-Client-Key: test" \
    -H "Content-Type: application/json" \
    -d '{"score":1337,"playerName":"Sean Plott","cheatMode":false}' \
    http://127.0.0.1:8080/v1/classes/GameScore
```

## 功能

## 开发日志

[开发日志.md](/开发日志.md)

## LICENSE

[MIT](/LICENSE)
