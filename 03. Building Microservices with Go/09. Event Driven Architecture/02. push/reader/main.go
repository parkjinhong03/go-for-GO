package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/nats-io/nats.go"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"runtime"
)

var address = flag.String("address", "localhost:4222", "NATS server URI")

type Product struct {
	Name string `json:"name" validate:"required"`
	SKU  string `json:"sku" validate:"required"`
}

func init() {
	flag.Parse()
}

func main() {
	natsCli, err := nats.Connect("nats://" + *address)
	if err != nil {
		log.Fatal(err)
	}

	_, err = natsCli.Subscribe("product.insert", MsgHandler)
	runtime.Goexit()
}

func MsgHandler(msg *nats.Msg) {
	var product Product
	decoder := json.NewDecoder(bytes.NewReader(msg.Data))
	if err := decoder.Decode(&product); err != nil {
		log.Printf("%s: Failed to decode event object | %v\n", msg.Subject, err)
		return
	}

	if err := validator.New().Struct(&product); err != nil {
		log.Printf("%s: Failed to decode event object | %v\n", msg.Subject, err)
		return
	}

	log.Printf("%s: Product(%s) was handled successfully\n", msg.Subject, product.Name)
}