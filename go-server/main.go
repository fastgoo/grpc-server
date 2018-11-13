package main

import (
	"golang.org/x/net/context"
	pb "local/sms/meta"
	"github.com/fastgoo/alisms-go"
	"github.com/Unknwon/goconfig"
	"log"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"path/filepath"
	"os"
)

const (
	port = ":9001"
)

type server struct{}

var alismsConfig map[string]string

func (s *server) Send(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	var defaultTemplate string
	if (in.Template != "") {
		defaultTemplate = in.Template;
	} else {
		defaultTemplate = alismsConfig["defaultTemplate"];
	}
	error := alisms.InitConfig("http://dysmsapi.aliyuncs.com", alismsConfig["accessKeyId"], alismsConfig["accessKeySecret"], alismsConfig["signName"]).Send(in.Mobile, in.Params, defaultTemplate)
	if error != nil {
		return nil, error
	}
	return &pb.Reply{Code: 200, Msg: "发送成功"}, nil
}

func main() {
	alismsConfig = loadConfig("alisms")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSmsServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//获取当前运行路径，用于取配置文件
func getDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

// 加载配置文件
// 这里如果是临时运行要注意写死地址，程序临时运行后会生成临时目录去运行，故而取不到config.ini
// 编译后的程序就是正常的
func loadConfig(field string) map[string]string {
	cfg, err := goconfig.LoadConfigFile(getDir() + "/config.ini")
	if err != nil {
		log.Println("读取配置文件失败[config.ini]")
	}
	value, _ := cfg.GetSection(field)
	return value
}
