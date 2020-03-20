// 이번 예제에서는 푸시 메시징 패턴에 대해 알아보고, 이를 구현해볼 것 이다.
// 푸시 패턴은 큐와 달리 이벤트가 발생하여 브로커에 등록되는 즉시 이 브로커를 구독하고 있는 모든 서비스둘이 메세지를 수신하는 패턴이다.
// 이 예제에서는 푸시 패턴을 NATS.io 라는 브로커를 이용하여 구현해 볼 것이다.

package main

import (
	"encoding/json"
	"flag"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
)

var address = flag.String("address", "", "NATS server URI")

type Product struct {
	Name string `json:"name"`
	SKU  string `json:"sku"`
}

func init() {
	flag.Parse()
}

type ProductHandler struct {
	NatsCli *nats.Conn
}

func (ph *ProductHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ph.insertProduct(rw, r)
	}
}

func (ph *ProductHandler) insertProduct(rw http.ResponseWriter, r *http.Request) {
	var product Product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&product); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(product)

	// *nat.Conn.Publish 메서드를 이용하여 원하는 제목과 페이로드를 가지고 있는 새로운 이벤트를 등록시킬 수 있다.
	if err := ph.NatsCli.Publish("product.insert", data); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("Succeeded to register product insert event to Nats")
}

func main() {
	natsCli, err := nats.Connect("nats://" + *address)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/product", &ProductHandler{NatsCli: natsCli})
	log.Println("Starting product write service on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}