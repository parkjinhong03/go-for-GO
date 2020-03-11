package main

import (
	"fmt"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"log"
)

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

	log.Fatal(service.Run())
}