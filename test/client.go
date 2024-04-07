package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var RedisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func Subscribe() {
	log.Println("Subscribing")
	pubSub := RedisClient.Subscribe(ctx, "mychannel")
	log.Println("Subscribed")
	for {
		log.Println("Receiving")
		msg, err := pubSub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Println(msg.Channel, msg.Payload)
	}
}

func main() {
	Subscribe()
}
