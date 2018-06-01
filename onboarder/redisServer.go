package main

import (
	redis "github.com/dotcloud/go-redis-server"
	"os"
	"log"
)

func main() {
	redisDB1()
}

func redisDB1() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatalf("Provide \"REDIS_ADDR\" environment variable")
	}
	server, err := redis.NewServer(redis.DefaultConfig())
	server.Addr = redisAddr
	if err != nil {
		panic(err)
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
