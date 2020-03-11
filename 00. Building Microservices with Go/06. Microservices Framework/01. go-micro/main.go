package main

import (
	"context"
	"fmt"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"go-micro.com/examples/1/proto"
	"log"
)

type Greeter struct {}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Message = "hello " + req.Name
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("greeter"),
		micro.Version("latest"),
		micro.Flags(
			cli.StringFlag{
				Name:  "env",
				Usage: "The Environment",
			},
		),
	)

	service.Init(
		micro.Action(func(ctx *cli.Context) {
			env := ctx.String("env")
			if len(env) > 0 {
				fmt.Println("Environment set to", env)
			}
		}),
	)

	err := proto.RegisterGreeterHandler(service.Server(), new(Greeter))
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(service.Run())
}