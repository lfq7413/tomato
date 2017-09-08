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

ServerURL = http://127.0.0.1:8080/v1
DatabaseType = MongoDB
DatabaseURI = 192.168.99.100:27017/test
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

## 启用 LiveQuery
###### 在 tomato 中添加配置项
```ini
LiveQueryClasses = classA|classB
PublisherType = Redis
PublisherURL = 192.168.99.100:6379
```
###### 启动 LiveQuery
使用默认参数与 tomato 同时启动：
```go
func main() {
    go tomato.RunLiveQueryServer(nil)
    tomato.Run()
}
```
使用自定义参数启动或者是独立运行：
```go
func main() {
    args := map[string]string{}
    args["pattern"] = "/livequery"
    args["addr"] = ":8089"
    args["logLevel"] = "VERBOSE"
    args["serverURL"] = "http://127.0.0.1:8080/v1"
    args["appId"] = "test"
    args["clientKey"] = "test"
    args["masterKey"] = "test"
    args["subType"] = "Redis"
    args["subURL"] = "192.168.99.100:6379"
    // 使用自定义参数与 tomato 同时启动
    go tomato.RunLiveQueryServer(args)
    tomato.Run()

    // 独立运行 LiveQueryServer
    // tomato.RunLiveQueryServer(args)
}
```

## 使用云代码
###### 使用云函数
声明：
```go
func main() {
    ...
    cloud.Define("hello", func(req cloud.FunctionRequest, resp cloud.Response) {
		// 函数的参数在 req.Params 中
		params := req.Params
		name := ""
		if params != nil && params["name"] != nil {
			if n, ok := params["name"].(string); ok {
				name = n
			}
		}
		// 计算结果通过 resp.Success() 返回
		resp.Success("hello " + name + "!")
		// 出现错误时通过 resp.Error() 返回错误
		// resp.Error(0, "error message")
	}, func(req cloud.FunctionRequest) bool {
		// 校验云函数是否可执行
		return true
	})
    ...
}
```
调用：
```bash
    curl -X POST \
    -H "X-Parse-Application-Id: test" \
    -H "X-Parse-Client-Key: test" \
    -H "Content-Type: application/json" \
    -d '{"name":"joe"}' \
    http://127.0.0.1:8080/v1/functions/hello
```
返回
```json
{
    "result": "hello joe!"
}
```
###### 使用 Hook 函数
BeforeSave
```go
func main() {
    ...
    cloud.BeforeSave("post", func(req cloud.TriggerRequest, resp cloud.Response) {
		// 在保存对象前可对对象进行直接修改
		if title, ok := req.Object["title"].(string); ok {
			if len(title) > 10 {
				req.Object["title"] = title[:10] + "..."
			}
		}
		// 当修改了 req.Object 时，可直接通过 resp.Success(nil) 返回成功
		// 当需要通过 resp.Success() 返回结果时，参数格式为 map[string]interface{}{"object": result}
		resp.Success(nil)
		// 出现错误时通过 resp.Error() 返回错误
		// resp.Error(0, "error message")
	})
    ...
}
```
AfterSave
```go
func main() {
    ...
    cloud.AfterSave("post", func(req cloud.TriggerRequest, resp cloud.Response) {
		// 在保存对象之后，可处理其他与该对象相关的事物
		fmt.Println("Save objectId", req.Object["objectId"])
		// 无需返回处理结果
	})
    ...
}
```
BeforeDelete
```go
func main() {
    ...
    cloud.BeforeDelete("post", func(req cloud.TriggerRequest, resp cloud.Response) {
		// 在删除对象前，可判断该对象是否可被删除
		fmt.Println("Delete objectId", req.Object["objectId"])
		// 可被删除则返回成功
		resp.Success(nil)
		// 不可删除则返回失败
		// resp.Error(0, "error message")
	})
    ...
}
```
AfterDelete
```go
func main() {
    ...
    cloud.AfterDelete("post", func(req cloud.TriggerRequest, resp cloud.Response) {
		// 在删除对象之后，可处理其他与该对象相关的事物
		fmt.Println("Delete objectId", req.Object["objectId"])
		// 无需返回处理结果
	})
    ...
}
```
BeforeFind
```go
func main() {
    ...
    cloud.BeforeFind("post", func(req cloud.TriggerRequest, resp cloud.Response) {
		// 在查询对象前，可对查询条件进行修改
		// 原始查询条件保存在 req.Query 中，可修改的字段包含 where limit skip include keys
		query := types.M{"limit": 5}
		// 修改之后的查询条件通过 resp.Success() 返回
		resp.Success(query)
		// 出现错误时通过 resp.Error() 返回错误
		// resp.Error(0, "error message")
	})
    ...
}
```
AfterFind
```go
func main() {
    ...
    cloud.AfterFind("post", func(req cloud.TriggerRequest, resp cloud.Response) {
		// 在查询对象后，可对查询结果进行修改
		// 原始查询结果保存在 req.Objects 中
		objects := req.Objects
		for _, o := range objects {
			if object := utils.M(o); object != nil {
				if title, ok := object["title"].(string); ok {
					object["title"] = title + " - by tomato"
				}
			}
		}
		// 修改之后的查询结果通过 resp.Success() 返回
		resp.Success(objects)
		// 出现错误时通过 resp.Error() 返回错误
		// resp.Error(0, "error message")
	})
    ...
}
```

## 功能

## 开发日志

[开发日志.md](/开发日志.md)

## LICENSE

[MIT](/LICENSE)
