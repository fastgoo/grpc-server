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
)

const (
	port = ":50051"
)

type server struct{}

var config map[string]string

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
	cfg, err := goconfig.LoadConfigFile("/Users/Mr.Zhou/Project/golang/grpc/go-server/config.ini")
	if err != nil {
		log.Println("读取配置文件失败[config.ini]")
	}
	value, _ := cfg.GetSection(field)
	config = value
}

func main() {
	loadConfig("alisms")
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
