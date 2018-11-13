### <p align="center">grpc 阿里云短信发送服务</p>
<p align="center">
  <a href="https://github.com/fastgoo/getui-php"><img src="https://img.shields.io/badge/license-MIT-brightgreen.svg"></a>
  <a href="https://golang.org/"><img src="https://img.shields.io/badge/golang->=10.0-brightgreen.svg"></a>
  <a href="https://php.net"><img src="https://img.shields.io/badge/grpc_server-go-brightgreen.svg"></a>
  <a href="https://php.net"><img src="https://img.shields.io/badge/grpc_client-php-brightgreen.svg"></a>
 <a href="https://php.net"><img src="https://img.shields.io/badge/php->=5.6-brightgreen.svg"></a>
  <a href="https://www.aliyun.com/product/sms"><img src="https://img.shields.io/badge/阿里云短信-grpc_服务-2077ff.svg"></a>
</p>

---

golang中文入门教程： https://golangcaff.com/docs/the-way-to-go

grpc-go服务端文档：https://grpc.io/docs/quickstart/go.html

protocol-buffers文档：https://developers.google.com/protocol-buffers/

php-grpc客户端文档：https://grpc.io/docs/quickstart/php.html

### 入门环境安装
1. 安装golang环境：https://golangcaff.com/docs/the-way-to-go/install-go-on-linux/8
2. 安装protoc环境：https://www.cnblogs.com/luoxn28/p/5303517.html
3. 安装grpc-go环境：https://www.jianshu.com/p/dba4c7a6d608
4. 安装php环境、grpc扩展、grpc_php_plugin插件：https://grpc.io/docs/quickstart/php.html


### grpc服务端

* **生成pb.go文件**

```bash
# --go_out 这里的地址是go类库地址，如果没有sms地址，请先创建地址才能正常生成
# 这条命令会在类库的sms文件夹中生成 meta/sms.pb.go 文件

protoc -I ./go-server/ --proto_path=./ --go_out=plugins=grpc:$GOPAHT/sms/ ./meta/sms.proto
```

* **修改main.go入口文件**

```golang
package main

import (
	"golang.org/x/net/context"
	pb "local/sms/meta" //这里是生成的pb文件所在的路径（需修改）
	"github.com/fastgoo/alisms-go"
	"github.com/Unknwon/goconfig"
	"log"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"//这里是设置端口
)

type server struct{}

var config map[string]string

//编写服务方法
func (s *server) Send(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	var defaultTemplate string
	if (in.Template != "") {
		defaultTemplate = in.Template;
	} else {
		defaultTemplate = config["defaultTemplate"];
	}
	error := alisms.InitConfig("http://dysmsapi.aliyuncs.com", config["accessKeyId"], config["accessKeySecret"], config["signName"]).Send(in.Mobile, in.Params, defaultTemplate)
	if error != nil {
		return nil, error
	}
	return &pb.Reply{Code: 200, Msg: "发送成功"}, nil
}

// 加载配置文件
func loadConfig(field string) {
   //这里的配置文件地址是我直接写死了我的真实地址，实际场景应该是使用相对地址（需修改）
	cfg, err := goconfig.LoadConfigFile("/Users/Mr.Zhou/Project/golang/grpc/go-server/config.ini")
	if err != nil {
		log.Println("读取配置文件失败[config.ini]")
	}
	value, _ := cfg.GetSection(field)
	config = value
}

func main() {
   //载入配置文件
	loadConfig("alisms")
	//监听tcp端口服务
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	//注册grpc服务
	pb.RegisterSmsServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

```


### grpc php客户端
* 生成php代码

```bash
# --proto_path 为输出路径
# --php_out 为生成PHP客户端代码所在的路径
# --grpc_out 为生成PHP依赖文件所在的路径
# --plugin 为PHP的插件路径（这个路径是grpc项目里面使用 make grpc_php_plugin 生成的路径）

protoc --proto_path=./ --php_out=./php-client/src/proto/ --grpc_out=./php-client/src/proto/ --plugin=protoc-gen-grpc=/Users/Mr.Zhou/Project/PHP/grpc-master/bins/opt/grpc_php_plugin ./meta/sms.proto
```

* php客户端代码

```php
require_once APP_DIR . '/vendor/autoload.php';

try {
    //初始化客户端
    $client = new \Sms\SmsClient("127.0.0.1:50051", [
        'credentials' => Grpc\ChannelCredentials::createInsecure()
    ]);
    //设置请求参数
    $request = new \Sms\Request();
    $request->setParams(json_encode(['code' => '123456']));
    $request->setMobile(15600087538);
    //$request->setTemplate("123");

    //获取返回值，这里的$error是一个对象数组，如果code不是0则发送失败了
    //这里$reply 为Reply的实例对象
    //wait为阻塞等待服务端返回参数
    list($reply, $error) = $client->Send($request)->wait();
    if ($error->code) {
        exit("发送失败，错误信息：" . $error->details);
    }
    var_dump($reply->getCode(), $reply->getMsg(), $reply->serializeToJsonString());

} catch (\Exception $exception) {
    var_dump($exception->getMessage());
}

```

* 安装composer依赖文件,执行客户端请求

```bash
# 安装compsoer 依赖库
composer install

# 执行客户端代码
php demo.php
```




