package common

// package for redis connection
// context is used to control request lifetimes, deadlines, and cancellations. Redis commands in Go require a context.
import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

// context is like a controller for redis commands
// context.Background gives us a base context - meaning there's no timeout or cancellation
// when you call commands later, eg. rdb.Ping(ctx), redis knows which context is managing the request.
var Ctx = context.Background()

func ConnectRedis() *redis.Client {

	// creating a redis client instance via redis.NewClient
	// redis.Options configures how to connect.
	// the variable rdb is now a redis client, and it'll be used for all Redis operations like GET, SET
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0, // default redis DB to use
	})

	// testing if the connection works
	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}

	log.Default().Println("Connected to redis successfully")
	return rdb

}
