package main

import (
	crane "github.com/artisticbones/Crane-RPC"
	"github.com/artisticbones/Crane-RPC/database"
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

func insert(key string, val string) error {
	err := database.RDB.Set(key, val, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func query(key string) (string, error) {
	val, err := database.RDB.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func main() {
	addr := "127.0.0.1:8080"
	srv := crane.NewServer(addr)
	database.RDB = database.InitRedis("redis://localhost:6379/")
	srv.Register("calcService", calcService)
	srv.Register("insert", insert)
	srv.Register("query", query)
	log.Println("service is running")
	go srv.Run()

	for {
	}
}
