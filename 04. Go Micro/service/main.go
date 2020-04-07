package main

import (
	"context"
	"fmt"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/parkjinhong03/go-micro-example/service/proto"
	"log"
	"time"
)

// RegisterGreeterHandler 함수에 등록하기 위한 객체 정의
type Greeter struct {}

// GreeterHandler 인터페이스로 사용하기 위해 Hello 메서드 정의 (실제로 rpc 요청이 들어왔을 때, 비즈니스 로직을 수행하는 메서드임)
func (g *Greeter) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}

// 클라이언트를 실행시키는 함수이다.
func RunClient(s micro.Service) {
	// protoc 명령어로 만든 proto.NewGreeterHandler 함수를 이용하여 greeter 서비스의 Greeter Handler에 대한 클라이언트를 생성할 수 있다.
	greeter := proto.NewGreeterService("greeter", s.Client())
	// greeter.Hello 메서드를 호출하여 원격 rpc 통신을 진행할 수 있다.
	// 참고로 context.WitTimeout 함수를 이용하면 기본 5초인 타임 아웃 시간을 원하는 시간으로 변경할 수 있다.
	resp, err := greeter.Hello(context.Background(), &proto.Request{Name: "Sample Client"})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Greeting)
}

func main() {
	// micro.NewService 함수를 이용하여 선택적으로 일정 옵션들을 포함한 새 서비스를 생성한다.
	service := micro.NewService(
		micro.Name("greeter"),
		micro.Version("latest"),
		micro.Metadata(map[string]string{
			"type": "helloworld",
		}),

		// micro.Flags 함수와 cli 패키지를 이용하여 서비스 실행 시 파싱할 flags 들을 설정할 수 있다.
		micro.Flags(&cli.BoolFlag{
			Name: "run_client",
			Usage: "Launch the client",
		}),
	)

	// protoc 명령어로 만든 proto.RegisterGreeterHandler 함수에 위에서 정의한 Greeter 핸들러를 넘겨 해당 핸들러를 서비스에 등록한다.
	if err := proto.RegisterGreeterHandler(service.Server(), new(Greeter)); err != nil {
		log.Fatal(err)
	}

	// service.Run 메서드를 이용해 해당 서비스를 서버로써 특정 포트에 바인딩하고, 서비스 레지스트리에 등록해준다.
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}