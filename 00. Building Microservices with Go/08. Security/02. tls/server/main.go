package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/plain")
		_, _ = fmt.Fprint(rw, "Hello, HTTPS!!")
	})

	log.Fatal(http.ListenAndServeTLS(":8080", "../key/instance_cert.pem", "../key/instance_key.pem", nil))
}