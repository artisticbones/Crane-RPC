package main

import (
	crane "github.com/artisticbones/Crane-RPC"
	"log"
	"net"
)

func main() {
	addr := "127.0.0.1:8080"
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("dial error: %v\n", err)
	}
	cli := crane.NewClient(conn)

	var callService func(string) (int, error)

	cli.Call("calcService", &callService)
	u, err := callService("abced")
	if err != nil {
		log.Printf("query error: %v\n", err)
	} else {
		log.Printf("query result: %v", u)
	}
	conn.Close()
}
