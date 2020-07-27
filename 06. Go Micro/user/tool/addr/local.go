package addr

import (
	"log"
	"net"
)

func GetLocal() *net.UDPAddr {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil { log.Fatal(err) }
	defer func() { _ = conn.Close() } ()
	return conn.LocalAddr().(*net.UDPAddr)
}