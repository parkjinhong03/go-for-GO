package main

import (
	"context"
	"fmt"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"go-micro.com/examples/1/client"
	"go-micro.com/examples/1/proto"
	"log"
	"os"
	"time"
)

type Greeter struct {}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Message = "Hello " + req.Name
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("greeter"),
		micro.Version("latest"),
		micro.Flags(
			cli.BoolFlag{
				Name:  "run_client",
				Usage: "Launch the client",
			},
			cli.StringFlag{
				Name:  "name",
				Usage: "Search Service's Name",
			},
		),
	)

	service.Init(
		micro.Action(func(ctx *cli.Context) {
			if !ctx.Bool("run_client") {
				return
			}
			if ctx.String("name") == "" {
				fmt.Println("Please specify service's name which is not blank if you want to run client.")
				return
			}
			go func() {
				time.Sleep(time.Second)
				resp := client.RunClient(service.Client(), ctx.String("name"), &proto.HelloRequest{
					Name: "Example",
				})
				fmt.Println(resp.Message)
				os.Exit(0)
			}()
		}),
	)

	err := proto.RegisterGreeterHandler(service.Server(), new(Greeter))
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(service.Run())
}