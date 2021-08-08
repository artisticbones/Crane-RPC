package main

import (
	crane "github.com/artisticbones/Crane-RPC"
	"log"
)

// calCalService 返回字符串中字符个数
func calcService(str string) (int, error) {
	sum := 0
	for _, v := range str {
		if v >= 'a' && v <= 'z' {
			sum++
		} else if v >= 'A' && v <= 'Z' {
			sum++
		}
	}

	return sum, nil
}

func main() {
	addr := "127.0.0.1:8080"
	srv := crane.NewServer(addr)
	srv.Register("calcService", calcService)
	log.Println("service is running")
	go srv.Run()

	for {
	}
}
