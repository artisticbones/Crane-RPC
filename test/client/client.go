package main

import (
	"fmt"
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
	var insertService func(string, string) error
	var queryService func(string) (string, error)

	cli.Call("calcService", &callService)
	cli.Call("insert", &insertService)
	cli.Call("query", &queryService)
	u, err := callService("abced")
	if err != nil {
		fmt.Printf("query error: %v\n", err)
	} else {
		fmt.Printf("query result: %v\n", u)
	}
	err = insertService("testKey", "testVal")
	if err != nil {
		fmt.Printf("insert error: %s\n", err.Error())
	} else {
		fmt.Printf("insert success!\n")
	}
	val, err := queryService("testKey")
	if err != nil {
		fmt.Printf("insert error: %s\n", err.Error())
	} else {
		fmt.Printf("query success! the val: %s\n", val)
	}
	err = insertService("wang-pi-dan", "1203")
	if err != nil {
		fmt.Printf("insert error: %s\n", err.Error())
	} else {
		fmt.Printf("insert success!\n")
	}
	val, err = queryService("wang-pi-dan")
	if err != nil {
		fmt.Printf("insert error: %s\n", err.Error())
	} else {
		fmt.Printf("query success! the val: %s\n", val)
	}
	conn.Close()
}
