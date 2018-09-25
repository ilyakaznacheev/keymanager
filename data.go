package main

import (
	"context"
	"log"

	"github.com/go-redis/redis"
)

type redisClient struct {
	client *redis.Client
	input  chan redisAction
	finish context.CancelFunc
}

type redisDo func(*redisClient)

type redisAction struct {
	action redisDo
}

var redClient *redisClient

// redisServer is used to execute change requests to Redis synchronously
// to prevent some bad things like race condition
func redisServer(ctx context.Context, client *redisClient) {
	for {
		select {
		case inputAction := <-client.input:
			inputAction.action(client)

		case <-ctx.Done():
			return
		}
	}
}

func initServer() {
	log.Println("server startup")

	config := getConfig()

	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		DB:       config.Redis.DataBase,
	})

	ctx, finish := context.WithCancel(context.Background())

	redClient = &redisClient{
		client: client,
		input:  make(chan redisAction),
		finish: finish,
	}

	go redisServer(ctx, redClient)
}

func shutDown() {
	redClient.finish()
	log.Println("server shutdown")
}
