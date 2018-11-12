# alisms go语言sdk

#### 项目介绍
基于go 实现的阿里云发送短信接口操作
我不是作者，我只是搬运工(来源)：https://blog.csdn.net/fyxichen/article/details/80896848


#### 安装教程

go get github.com/fastgoo/alisms-go

#### 使用说明

```golang
package main

import (
    "fmt"
	"github.com/fastgoo/alisms-go"
)

func main()
{
    error := alisms.InitConfig("http://dysmsapi.aliyuncs.com", "accessKeyId", "accessKeySecret", "signName").Send("手机", `{"code":"123"}`, "模板编号")
    if error != nil {
    	//发送失败，获取错误消息
    	fmt.Print(error)
    }else{
        fmt.Print("发送成功")
    }
}

```