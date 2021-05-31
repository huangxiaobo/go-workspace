package utils

import (
	"fmt"
	"log"
)

var (
	Client *redis.Client
)

func Connect() (*redis.Client, error) {
	if Client != nil {
		return Client, nil
	}

	option := redis.Options{
		Addr: "0.0.0.0:6379",
		//Password: "123456", // no password set
		DB: 0, // use default DB
	}
	Client = redis.NewClient(&option)

	pong, err := Client.Ping().Result()
	if err != nil {
		fmt.Println("redis 连接失败：", pong, err)
		return nil, err
	}
	fmt.Println("redis 连接成功：", pong)
	return Client, nil
}

func Close() {
	err := Client.Close()
	if err != nil {
		log.Fatal(err)
	}
}
